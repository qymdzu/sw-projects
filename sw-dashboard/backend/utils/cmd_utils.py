"""
系统命令工具 — 异步执行 shell 命令、systemctl 操作、crontab 管理
"""

import asyncio
import shlex
from typing import Dict, List, Optional


class CommandTimeoutError(Exception):
    """命令执行超时异常"""
    pass


async def run_cmd(command: str, timeout: int = 10) -> Dict:
    """
    异步执行 shell 命令，返回结果。

    Args:
        command: shell 命令字符串
        timeout: 超时秒数

    Returns:
        {"returncode": int, "stdout": str, "stderr": str}
    """
    try:
        proc = await asyncio.create_subprocess_shell(
            command,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        try:
            stdout, stderr = await asyncio.wait_for(
                proc.communicate(), timeout=timeout
            )
            return {
                "returncode": proc.returncode or 0,
                "stdout": stdout.decode("utf-8", errors="replace"),
                "stderr": stderr.decode("utf-8", errors="replace"),
            }
        except asyncio.TimeoutError:
            try:
                proc.kill()
                await proc.wait()
            except Exception:
                pass
            raise CommandTimeoutError(
                f"Command timed out after {timeout}s: {command[:100]}"
            )
    except FileNotFoundError:
        return {"returncode": -1, "stdout": "", "stderr": "command not found"}
    except PermissionError:
        return {"returncode": -1, "stdout": "", "stderr": "permission denied"}


async def run_systemctl(action: str, service_name: str = "hermes-gateway") -> Dict:
    """
    执行 systemctl 操作。

    Args:
        action: "is-active" | "status" | "show"
        service_name: 服务名

    Returns:
        解析后的状态字典
    """
    try:
        result = await run_cmd(f"systemctl {action} {service_name}", timeout=5)
    except CommandTimeoutError:
        return {"status": "unknown"}

    stdout = result.get("stdout", "").strip()
    returncode = result.get("returncode", -1)

    if action == "is-active":
        if returncode == 0:
            return {"status": stdout or "active"}
        return {"status": "inactive"}

    if action == "status":
        if returncode != 0:
            return {"status": "inactive"}
        info = {"status": "active"}
        for line in stdout.split("\n"):
            line = line.strip()
            if "Main PID:" in line:
                try:
                    info["pid"] = int(line.split("Main PID:")[1].split()[0])
                except (ValueError, IndexError):
                    pass
            if "Memory:" in line:
                try:
                    mem_str = line.split("Memory:")[1].strip().split()[0]
                    if mem_str.endswith("M"):
                        info["memory_mb"] = float(mem_str[:-1])
                    elif mem_str.endswith("G"):
                        info["memory_mb"] = float(mem_str[:-1]) * 1024
                    elif mem_str.endswith("K"):
                        info["memory_mb"] = float(mem_str[:-1]) / 1024
                    else:
                        try:
                            info["memory_mb"] = float(mem_str) / (1024 * 1024)
                        except ValueError:
                            pass
                except (ValueError, IndexError):
                    pass
            if "CPUUsageNSec=" in line:
                try:
                    ns = int(line.split("=")[1].strip())
                    info["cpu_percent"] = round(ns / 1e7, 1)  # rough estimate
                except (ValueError, IndexError):
                    pass
            if "uptime" in line.lower() or "Active:" in line:
                if "ago" in line:
                    info["uptime"] = line.split(";")[-1].strip() if ";" in line else line.strip()
        return info

    if action == "show":
        info = {"status": "active"}
        for line in stdout.split("\n"):
            line = line.strip()
            if "=" in line:
                key, _, val = line.partition("=")
                key = key.strip()
                val = val.strip()
                if key == "MainPID" and val and val != "0":
                    try:
                        info["pid"] = int(val)
                    except ValueError:
                        pass
                if key == "MemoryCurrent":
                    try:
                        info["memory_mb"] = round(int(val) / (1024 * 1024), 2)
                    except (ValueError, IndexError):
                        pass
                if key == "CPUUsageNSec":
                    try:
                        info["cpu_percent"] = round(int(val) / 1e9, 2)
                    except (ValueError, IndexError):
                        pass
        return info

    return {"status": "unknown"}


async def run_crontab(action: str, command: str = None, enable: bool = True) -> Dict:
    """
    管理 crontab 任务。

    Args:
        action: "list" | "toggle"
        command: toggle 模式下要匹配的命令
        enable: toggle 模式下 True=启用, False=禁用

    Returns:
        list: {"jobs": [...], "success": True}
        toggle: {"success": bool, "message": str}
    """
    if action == "list":
        try:
            result = await run_cmd("crontab -l 2>/dev/null || true", timeout=5)
        except CommandTimeoutError:
            return {"jobs": [], "success": True}

        stdout = result.get("stdout", "")
        jobs = []
        for line in stdout.split("\n"):
            line = line.strip()
            if not line:
                continue
            # Pure comment line (only #)
            if line.startswith("#") and not _looks_like_cron(line.lstrip("#").strip()):
                continue

            enabled = not line.startswith("#")
            raw_line = line
            if not enabled:
                raw_line = line.lstrip("#").strip()

            # Parse cron expression (5 fields + command)
            parts = raw_line.split()
            if len(parts) < 6:
                continue  # Not a valid cron line

            cron_parts = parts[:5]
            cmd = " ".join(parts[5:])
            schedule = " ".join(cron_parts)

            comment = ""
            comment_idx = cmd.find("#")
            if comment_idx >= 0:
                comment = cmd[comment_idx + 1 :].strip()
                cmd = cmd[:comment_idx].strip()

            jobs.append({
                "raw": line,
                "schedule": schedule,
                "command": cmd,
                "comment": comment,
                "enabled": enabled,
                "last_run": None,
                "last_output": None,
            })

        return {"jobs": jobs, "success": True}

    if action == "toggle":
        if not command:
            return {"success": False, "message": "command 参数不能为空"}
        try:
            list_result = await run_cmd("crontab -l 2>/dev/null || true", timeout=5)
        except CommandTimeoutError:
            return {"success": False, "message": "crontab -l 执行超时"}

        lines = list_result.get("stdout", "").split("\n")
        new_lines = []
        matched = False

        for line in lines:
            stripped = line.strip()
            # Check both enabled and disabled forms
            if command in stripped or command in stripped.lstrip("#"):
                matched = True
                if enable:
                    # Remove leading # and whitespace
                    new_line = stripped.lstrip("#").strip()
                else:
                    # Add # at start
                    new_line = "# " + stripped if not stripped.startswith("#") else stripped
                new_lines.append(new_line)
            else:
                new_lines.append(line)

        if not matched:
            return {"success": False, "message": "未找到匹配的 cron 任务"}

        content = "\n".join(new_lines).strip() + "\n"
        try:
            # Use python to write to crontab (avoids shell escaping issues)
            import sys
            write_result = await run_cmd(
                f'python3 -c "import sys; sys.stdout.write({shlex.quote(content)})" | crontab -',
                timeout=5,
            )
        except CommandTimeoutError:
            return {"success": False, "message": "写入 crontab 超时"}

        if write_result.get("returncode", -1) == 0:
            action_text = "已启用" if enable else "已禁用"
            return {"success": True, "message": f"Cron 任务 {action_text}"}
        else:
            return {
                "success": False,
                "message": f"写入 crontab 失败: {write_result.get('stderr', '')}",
            }

    return {"success": False, "message": f"未知操作: {action}"}


def _looks_like_cron(text: str) -> bool:
    """判断文本是否看起来像 cron 表达式"""
    parts = text.split()
    if len(parts) < 6:
        return False
    # Check first 5 parts are cron-like
    cron_chars = set("0123456789,/*-")
    for p in parts[:5]:
        for c in p:
            if c not in cron_chars:
                return False
    return True
