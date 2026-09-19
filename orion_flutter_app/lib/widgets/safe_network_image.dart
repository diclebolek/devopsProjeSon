import 'package:flutter/material.dart';
import '../constants/colors.dart';

/// Ağ görseli başarısız olunca ImageCodecException / overflow göstermez.
Widget safeNetworkImage(
  String? url, {
  double? width,
  double? height,
  BoxFit fit = BoxFit.cover,
  IconData fallbackIcon = Icons.fitness_center,
  Color? fallbackColor,
}) {
  final iconColor = fallbackColor ?? SiriusColors.accent;
  final placeholder = Container(
    width: width,
    height: height,
    color: iconColor.withValues(alpha: 0.12),
    alignment: Alignment.center,
    child: Icon(fallbackIcon, color: iconColor, size: (height ?? 40) * 0.45),
  );

  if (url == null || url.trim().isEmpty || !url.startsWith('http')) {
    return placeholder;
  }

  return Image.network(
    url,
    width: width,
    height: height,
    fit: fit,
    errorBuilder: (_, __, ___) => placeholder,
    loadingBuilder: (context, child, progress) {
      if (progress == null) return child;
      return placeholder;
    },
  );
}

/// CircleAvatar yerine güvenli avatar (NetworkImage kullanmaz).
Widget safeAvatar({
  String? imageUrl,
  double radius = 24,
  String? initials,
  IconData icon = Icons.person,
}) {
  final size = radius * 2;
  return ClipOval(
    child: SizedBox(
      width: size,
      height: size,
      child: safeNetworkImage(
        imageUrl,
        width: size,
        height: size,
        fallbackIcon: icon,
      ),
    ),
  );
}
