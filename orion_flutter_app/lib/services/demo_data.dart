import '../models/appointment.dart';
import '../models/customer.dart';
import '../models/employee.dart';
import '../models/service.dart';

/// Ağ/Supabase beklemeden UI'ı dolduran lokal demo veriler.
class DemoData {
  static const String demoEmail = 'test@gmail.com';
  static const String demoPhone = '0555 123 45 67';
  static const String demoName = 'Test';
  static const String demoLastName = 'Kullanıcı';

  static String image(String seed, {int w = 800, int h = 600}) =>
      'https://picsum.photos/seed/$seed/$w/$h';

  static Customer demoCustomer({String? email}) {
    final e = (email == null || email.trim().isEmpty) ? demoEmail : email.trim();
    final local = e.contains('@') ? e.split('@').first : e;
    return Customer(
      customerId: 1001,
      firstName: local.isEmpty ? demoName : _title(local),
      lastName: demoLastName,
      email: e,
      phone: demoPhone,
      address: 'Atatürk Cad. No:12, İstanbul',
      birthDate: DateTime(1995, 5, 15),
      gender: 'Diğer',
      createdAt: DateTime.now().subtract(const Duration(days: 120)),
      isActive: true,
    );
  }

  static Employee demoAdmin({String? email}) {
    return Employee(
      id: 1,
      firstName: 'Admin',
      lastName: 'Orion',
      expertise: 'Yönetim',
      skills: 'admin',
      email: email ?? 'admin@orion.com',
      phone: '0555 000 00 01',
      hireDate: DateTime(2020, 1, 1),
      isActive: true,
      profileImage: image('admin-avatar', w: 200, h: 200),
    );
  }

  static List<Service> services() => [
        Service(
          serviceId: 1,
          serviceName: 'Kişisel Antrenman',
          serviceDuration: 60,
          servicePrice: 150,
          description: 'Birebir fitness seansı',
          imageUrl: image('svc-pt'),
          category: 'Fitness',
          isActive: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        Service(
          serviceId: 2,
          serviceName: 'Grup Fitness',
          serviceDuration: 45,
          servicePrice: 80,
          description: 'Enerjik grup dersi',
          imageUrl: image('svc-group'),
          category: 'Fitness',
          isActive: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        Service(
          serviceId: 3,
          serviceName: 'Yoga',
          serviceDuration: 75,
          servicePrice: 100,
          description: 'Esneklik ve denge',
          imageUrl: image('svc-yoga'),
          category: 'Wellness',
          isActive: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        Service(
          serviceId: 4,
          serviceName: 'Pilates',
          serviceDuration: 50,
          servicePrice: 120,
          description: 'Core güçlendirme',
          imageUrl: image('svc-pilates'),
          category: 'Wellness',
          isActive: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      ];

  static List<Employee> employees() => [
        Employee(
          id: 10,
          firstName: 'Zeynep',
          lastName: 'Trainer',
          expertise: 'Fitness',
          skills: 'PT, HIIT',
          email: 'zeynep@orion.com',
          phone: '0555 111 22 33',
          hireDate: DateTime(2022, 3, 1),
          isActive: true,
          profileImage: image('emp-zeynep', w: 400, h: 400),
        ),
        Employee(
          id: 11,
          firstName: 'Can',
          lastName: 'Coach',
          expertise: 'Yoga',
          skills: 'Yoga, Pilates',
          email: 'can@orion.com',
          phone: '0555 222 33 44',
          hireDate: DateTime(2021, 6, 15),
          isActive: true,
          profileImage: image('emp-can', w: 400, h: 400),
        ),
        Employee(
          id: 12,
          firstName: 'Elif',
          lastName: 'Demir',
          expertise: 'Pilates',
          skills: 'Pilates, Stretch',
          email: 'elif@orion.com',
          phone: '0555 333 44 55',
          hireDate: DateTime(2023, 1, 10),
          isActive: true,
          profileImage: image('emp-elif', w: 400, h: 400),
        ),
      ];

  /// Bugün 12:00 dahil görünür test randevuları.
  static List<Appointment> appointments({String? forEmail}) {
    final email = forEmail ?? demoEmail;
    final customer = demoCustomer(email: email);
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final emps = employees();
    final svcs = services();

    Appointment build({
      required int id,
      required DateTime when,
      required String status,
      required int empIdx,
      required int svcIdx,
      String process = 'Waiting',
      String notes = '',
    }) {
      final emp = emps[empIdx % emps.length];
      final svc = svcs[svcIdx % svcs.length];
      return Appointment(
        appointmentId: id,
        randevuId: 'demo-$id',
        calisanId: emp.id,
        customerName: customer.fullName,
        employeeName: emp.fullName,
        serviceName: svc.serviceName,
        process: process,
        totalPrice: svc.servicePrice,
        appointmentDateTime: when,
        approvalStatus: status,
        notes: notes.isEmpty ? 'Demo randevu' : notes,
        customerPhone: customer.phone,
        customerEmail: customer.email,
        createdAt: now.subtract(const Duration(days: 2)),
        updatedAt: now,
      );
    }

    return [
      // Bugün 12:00 — en görünür
      build(
        id: 1,
        when: today.add(const Duration(hours: 12)),
        status: 'Approved',
        empIdx: 0,
        svcIdx: 0,
        process: 'Waiting',
        notes: 'Bugün 12:00 test randevusu',
      ),
      build(
        id: 2,
        when: today.add(const Duration(hours: 15, minutes: 30)),
        status: 'Pending',
        empIdx: 1,
        svcIdx: 2,
        notes: 'Bugün 15:30',
      ),
      build(
        id: 3,
        when: today.add(const Duration(days: 1, hours: 12)),
        status: 'Approved',
        empIdx: 2,
        svcIdx: 3,
        notes: 'Yarın 12:00',
      ),
      build(
        id: 4,
        when: today.add(const Duration(days: 2, hours: 10)),
        status: 'Pending',
        empIdx: 0,
        svcIdx: 1,
      ),
      build(
        id: 5,
        when: today.subtract(const Duration(days: 1)).add(const Duration(hours: 12)),
        status: 'Approved',
        empIdx: 1,
        svcIdx: 0,
        process: 'Completed',
        notes: 'Dün 12:00 tamamlandı',
      ),
      build(
        id: 6,
        when: today.subtract(const Duration(days: 3)).add(const Duration(hours: 14)),
        status: 'Canceled',
        empIdx: 2,
        svcIdx: 2,
        process: 'Completed',
      ),
      // Diğer müşteri (admin listesinde çeşitlilik)
      Appointment(
        appointmentId: 7,
        randevuId: 'demo-7',
        calisanId: 10,
        customerName: 'Ayşe Yılmaz',
        employeeName: emps[0].fullName,
        serviceName: svcs[1].serviceName,
        process: 'Waiting',
        totalPrice: svcs[1].servicePrice,
        appointmentDateTime: today.add(const Duration(hours: 13)),
        approvalStatus: 'Pending',
        notes: 'Diğer müşteri',
        customerPhone: '0555 987 65 43',
        customerEmail: 'ayse@orion.com',
        createdAt: now,
        updatedAt: now,
      ),
      Appointment(
        appointmentId: 8,
        randevuId: 'demo-8',
        calisanId: 11,
        customerName: 'Mehmet Demir',
        employeeName: emps[1].fullName,
        serviceName: svcs[2].serviceName,
        process: 'In Progress',
        totalPrice: svcs[2].servicePrice,
        appointmentDateTime: today.add(const Duration(days: 3, hours: 11)),
        approvalStatus: 'Approved',
        notes: 'Admin görünümü için',
        customerPhone: '0555 444 55 66',
        customerEmail: 'mehmet@orion.com',
        createdAt: now,
        updatedAt: now,
      ),
    ];
  }

  static List<Appointment> appointmentsForCustomer(String email) {
    final target = email.trim().isEmpty ? demoEmail : email.trim().toLowerCase();
    return appointments(forEmail: target)
        .where((a) => (a.customerEmail ?? '').toLowerCase() == target)
        .toList();
  }

  /// Profil: müşteri randevuları; yoksa demo müşteri listesi (UI boş kalmasın).
  static List<Appointment> profileAppointments(String email) {
    final mine = appointmentsForCustomer(email);
    if (mine.isNotEmpty) return mine;
    return appointmentsForCustomer(demoEmail);
  }

  static Map<String, dynamic> isletme() => {
        'isim': 'Orion Gym & Fitness',
        'aciklama': 'Spor Salonu & Fitness',
        'banner_url': image('orion-banner', w: 1200, h: 600),
        'logo_url': image('orion-logo', w: 200, h: 200),
        'arka_plan_url': image('orion-bg', w: 1200, h: 800),
        'telefon': '0212 555 00 00',
        'email': 'info@orion.com',
        'adres': 'Caddebostan, İstanbul',
      };

  static List<String> gallery() => [
        image('gallery-1'),
        image('gallery-2'),
        image('gallery-3'),
        image('gallery-4'),
        image('gallery-5'),
        image('gallery-6'),
      ];

  static List<String> events() => [
        image('event-1', w: 900, h: 500),
        image('event-2', w: 900, h: 500),
        image('event-3', w: 900, h: 500),
      ];

  static List<Map<String, String>> team() => [
        {
          'name': 'Zeynep Trainer',
          'role': 'Fitness Uzmanı',
          'image': image('team-1', w: 400, h: 500),
        },
        {
          'name': 'Can Coach',
          'role': 'Yoga Eğitmeni',
          'image': image('team-2', w: 400, h: 500),
        },
        {
          'name': 'Elif Demir',
          'role': 'Pilates Uzmanı',
          'image': image('team-3', w: 400, h: 500),
        },
      ];

  static Map<String, String> icerikBlok() => {
        'why_1_title': 'Güçlü Başla',
        'why_1_desc': 'Kişisel hedeflerin için profesyonel destek',
        'why_2_title': 'Esnek Program',
        'why_2_desc': 'Güne ve seviyene uygun dersler',
        'why_3_title': 'Topluluk',
        'why_3_desc': 'Motivasyon dolu bir spor ortamı',
        'currency': 'TL',
        'dialog_close': 'Kapat',
        'hero_subtitle': 'Spor Salonu & Fitness',
        'hero_tagline': 'Güçlü Hisset, Orion Ol',
        'brand_name': 'Orion',
        'btn_login': 'Giriş Yap',
        'btn_learn_more': 'Daha Fazla',
      };

  static String _title(String s) {
    if (s.isEmpty) return demoName;
    return s[0].toUpperCase() + s.substring(1);
  }
}
