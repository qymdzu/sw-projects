// lib/core/data/models/user.dart
// 用户模型
class User {
  final String id;
  final String name;
  final String role;  // 'student' / 'adult'
  final String? grade;  // 'primary_1' ... 'senior_12'
  final DateTime createdAt;
  final DateTime? lastActiveAt;

  const User({
    required this.id,
    required this.name,
    required this.role,
    this.grade,
    required this.createdAt,
    this.lastActiveAt,
  });

  bool get isStudent => role == 'student';
  bool get isAdult => role == 'adult';

  bool get isHighSchool => grade?.startsWith('senior_') ?? false;
  bool get isJunior => grade?.startsWith('junior_') ?? false;
  bool get isPrimary => grade?.startsWith('primary_') ?? false;

  Map<String, Object?> toMap() => {
    'id': id,
    'name': name,
    'role': role,
    'grade': grade,
    'created_at': createdAt.millisecondsSinceEpoch,
    'last_active_at': lastActiveAt?.millisecondsSinceEpoch,
  };

  factory User.fromMap(Map<String, Object?> m) => User(
    id: m['id']! as String,
    name: m['name']! as String,
    role: m['role']! as String,
    grade: m['grade'] as String?,
    createdAt: DateTime.fromMillisecondsSinceEpoch(m['created_at']! as int),
    lastActiveAt: m['last_active_at'] != null
        ? DateTime.fromMillisecondsSinceEpoch(m['last_active_at']! as int)
        : null,
  );

  User copyWith({String? name, String? grade, DateTime? lastActiveAt}) => User(
    id: id,
    name: name ?? this.name,
    role: role,
    grade: grade ?? this.grade,
    createdAt: createdAt,
    lastActiveAt: lastActiveAt ?? this.lastActiveAt,
  );
}
