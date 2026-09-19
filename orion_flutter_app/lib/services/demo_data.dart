import '../models/appointment.dart';
import '../models/customer.dart';
import '../models/employee.dart';
import '../models/service.dart';

/// Spor salonu (gym) odaklı lokal demo — boş ekran bırakmaz.
class DemoData {
  static const String demoEmail = 'test@gmail.com';
  static const String demoPhone = '0555 123 45 67';
  static const String demoName = 'Test';
  static const String demoLastName = 'Kullanıcı';

  /// Unsplash fitness görselleri (sabit, boşluk bırakmaz)
  static String image(String seed, {int w = 800, int h = 600}) {
    const map = <String, String>{
      'orion-banner':
          'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?auto=format&fit=crop&w=720&h=400&q=60',
      'orion-logo':
          'https://images.unsplash.com/photo-1571902943202-507ec2618e8f?auto=format&fit=crop&w=200&h=200&q=60',
      'orion-bg':
          'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?auto=format&fit=crop&w=720&h=480&q=60',
      'svc-pt':
          'https://images.unsplash.com/photo-1571019614242-c5c5dee9f50b?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-group':
          'https://images.unsplash.com/photo-1518611012118-696072aa579a?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-yoga':
          'https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-pilates':
          'https://images.unsplash.com/photo-1518611012118-696072aa579a?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-hiit':
          'https://images.unsplash.com/photo-1599058945522-28d584b6f14f?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-spin':
          'https://images.unsplash.com/photo-1534367507873-d2d7e24c797f?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-crossfit':
          'https://images.unsplash.com/photo-1526506118085-60ce8714f8c5?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-boxing':
          'https://images.unsplash.com/photo-1549719386-9dddf0c3b749?auto=format&fit=crop&w=480&h=360&q=60',
      'svc-swim':
          'https://images.unsplash.com/photo-1519315901367-f34ff9119d6f?auto=format&fit=crop&w=480&h=360&q=60',
      'emp-zeynep':
          'https://images.unsplash.com/photo-1594381898411-846e7d193883?auto=format&fit=crop&w=200&h=200&q=60',
      'emp-can':
          'https://images.unsplash.com/photo-1567013127542-490d757e51fc?auto=format&fit=crop&w=200&h=200&q=60',
      'emp-elif':
          'https://images.unsplash.com/photo-1518310383802-640c2de311b2?auto=format&fit=crop&w=200&h=200&q=60',
      'emp-burak':
          'https://images.unsplash.com/photo-1583454110551-21f2fa2afe61?auto=format&fit=crop&w=200&h=200&q=60',
      'emp-selin':
          'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?auto=format&fit=crop&w=200&h=200&q=60',
      'emp-mert':
          'https://images.unsplash.com/photo-1581009146145-b5ef439e1539?auto=format&fit=crop&w=200&h=200&q=60',
      'team-1':
          'https://images.unsplash.com/photo-1594381898411-846e7d193883?auto=format&fit=crop&w=280&h=350&q=60',
      'team-2':
          'https://images.unsplash.com/photo-1567013127542-490d757e51fc?auto=format&fit=crop&w=280&h=350&q=60',
      'team-3':
          'https://images.unsplash.com/photo-1518310383802-640c2de311b2?auto=format&fit=crop&w=280&h=350&q=60',
      'team-4':
          'https://images.unsplash.com/photo-1583454110551-21f2fa2afe61?auto=format&fit=crop&w=280&h=350&q=60',
      'team-5':
          'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?auto=format&fit=crop&w=280&h=350&q=60',
      'team-6':
          'https://images.unsplash.com/photo-1581009146145-b5ef439e1539?auto=format&fit=crop&w=280&h=350&q=60',
      'gallery-1':
          'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?auto=format&fit=crop&w=480&h=360&q=60',
      'gallery-2':
          'https://images.unsplash.com/photo-1571902943202-507ec2618e8f?auto=format&fit=crop&w=480&h=360&q=60',
      'gallery-3':
          'https://images.unsplash.com/photo-1540497077202-7c8a3999166f?auto=format&fit=crop&w=480&h=360&q=60',
      'gallery-4':
          'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?auto=format&fit=crop&w=480&h=360&q=60',
      'gallery-5':
          'https://images.unsplash.com/photo-1558611848-73f7eb4001a1?auto=format&fit=crop&w=480&h=360&q=60',
      'gallery-6':
          'https://images.unsplash.com/photo-1576678927484-cc907957088c?auto=format&fit=crop&w=480&h=360&q=60',
      'event-1':
          'https://images.unsplash.com/photo-1476480862126-209bfaa8edc8?auto=format&fit=crop&w=640&h=360&q=60',
      'event-2':
          'https://images.unsplash.com/photo-1434682881908-b43d0467b798?auto=format&fit=crop&w=640&h=360&q=60',
      'event-3':
          'https://images.unsplash.com/photo-1517838277536-f5f99be501cd?auto=format&fit=crop&w=640&h=360&q=60',
      'profile-avatar':
          'https://images.unsplash.com/photo-1633332755192-727a05c4013d?auto=format&fit=crop&w=300&h=300&q=60',
      'admin-avatar':
          'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=200&h=200&q=60',
    };
    return map[seed] ??
        'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?auto=format&fit=crop&w=$w&h=$h&q=60';
  }

  static Customer demoCustomer({String? email}) {
    final e =
        (email == null || email.trim().isEmpty) ? demoEmail : email.trim();
    final local = e.contains('@') ? e.split('@').first : e;
    return Customer(
      customerId: 1001,
      firstName: local.isEmpty ? demoName : _title(local),
      lastName: demoLastName,
      email: e,
      phone: demoPhone,
      address: 'Caddebostan Mah. Spor Cad. No:12, İstanbul',
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
      profileImage: image('admin-avatar'),
    );
  }

  static List<Service> services() {
    final now = DateTime.now();
    return [
      Service(
        serviceId: 1,
        serviceName: 'Kişisel Antrenman',
        serviceDuration: 60,
        servicePrice: 450,
        description: 'Birebir PT — hedefe özel program',
        imageUrl: image('svc-pt'),
        category: 'Fitness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 2,
        serviceName: 'Grup Fitness',
        serviceDuration: 45,
        servicePrice: 120,
        description: 'Enerjik grup dersi, her seviyeye uygun',
        imageUrl: image('svc-group'),
        category: 'Fitness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 3,
        serviceName: 'HIIT',
        serviceDuration: 40,
        servicePrice: 150,
        description: 'Yüksek yoğunluklu interval antrenman',
        imageUrl: image('svc-hiit'),
        category: 'Fitness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 4,
        serviceName: 'CrossFit',
        serviceDuration: 55,
        servicePrice: 180,
        description: 'Fonksiyonel güç ve kondisyon',
        imageUrl: image('svc-crossfit'),
        category: 'Fitness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 5,
        serviceName: 'Yoga',
        serviceDuration: 75,
        servicePrice: 140,
        description: 'Esneklik, denge ve nefes çalışması',
        imageUrl: image('svc-yoga'),
        category: 'Wellness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 6,
        serviceName: 'Pilates',
        serviceDuration: 50,
        servicePrice: 160,
        description: 'Core güçlendirme ve postür',
        imageUrl: image('svc-pilates'),
        category: 'Wellness',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 7,
        serviceName: 'Spinning',
        serviceDuration: 45,
        servicePrice: 130,
        description: 'Kapalı bisiklet — kardiyo odaklı',
        imageUrl: image('svc-spin'),
        category: 'Kardiyo',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
      Service(
        serviceId: 8,
        serviceName: 'Boks',
        serviceDuration: 50,
        servicePrice: 170,
        description: 'Teknik + kondisyon boks seansı',
        imageUrl: image('svc-boxing'),
        category: 'Spor',
        isActive: true,
        createdAt: now,
        updatedAt: now,
      ),
    ];
  }

  static List<Employee> employees() => [
        Employee(
          id: 10,
          firstName: 'Zeynep',
          lastName: 'Yılmaz',
          expertise: 'Kişisel Antrenman',
          skills: 'PT, HIIT, Beslenme',
          email: 'zeynep@orion.com',
          phone: '0555 111 22 33',
          hireDate: DateTime(2022, 3, 1),
          isActive: true,
          profileImage: null,
          prolificacy: 4.8,
        ),
        Employee(
          id: 11,
          firstName: 'Can',
          lastName: 'Öztürk',
          expertise: 'Yoga',
          skills: 'Yoga, Stretching',
          email: 'can@orion.com',
          phone: '0555 222 33 44',
          hireDate: DateTime(2021, 6, 15),
          isActive: true,
          profileImage: null,
          prolificacy: 4.7,
        ),
        Employee(
          id: 12,
          firstName: 'Elif',
          lastName: 'Demir',
          expertise: 'Pilates',
          skills: 'Pilates, Reformer',
          email: 'elif@orion.com',
          phone: '0555 333 44 55',
          hireDate: DateTime(2023, 1, 10),
          isActive: true,
          profileImage: null,
          prolificacy: 4.9,
        ),
        Employee(
          id: 13,
          firstName: 'Burak',
          lastName: 'Kaya',
          expertise: 'CrossFit',
          skills: 'CrossFit, Güç',
          email: 'burak@orion.com',
          phone: '0555 444 55 66',
          hireDate: DateTime(2020, 9, 1),
          isActive: true,
          profileImage: null,
          prolificacy: 4.6,
        ),
        Employee(
          id: 14,
          firstName: 'Selin',
          lastName: 'Arslan',
          expertise: 'Grup Fitness',
          skills: 'Zumba, Aerobik',
          email: 'selin@orion.com',
          phone: '0555 555 66 77',
          hireDate: DateTime(2022, 11, 20),
          isActive: true,
          profileImage: null,
          prolificacy: 4.8,
        ),
        Employee(
          id: 15,
          firstName: 'Mert',
          lastName: 'Çelik',
          expertise: 'Boks',
          skills: 'Boks, Kickboks',
          email: 'mert@orion.com',
          phone: '0555 666 77 88',
          hireDate: DateTime(2021, 4, 5),
          isActive: true,
          profileImage: null,
          prolificacy: 4.5,
        ),
      ];

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
      build(
        id: 1,
        when: today.add(const Duration(hours: 12)),
        status: 'Approved',
        empIdx: 0,
        svcIdx: 0,
        notes: 'Bugün 12:00 — Kişisel Antrenman',
      ),
      build(
        id: 2,
        when: today.add(const Duration(hours: 15, minutes: 30)),
        status: 'Pending',
        empIdx: 1,
        svcIdx: 4,
        notes: 'Bugün 15:30 — Yoga',
      ),
      build(
        id: 3,
        when: today.add(const Duration(days: 1, hours: 12)),
        status: 'Approved',
        empIdx: 2,
        svcIdx: 5,
        notes: 'Yarın 12:00 — Pilates',
      ),
      build(
        id: 4,
        when: today.add(const Duration(days: 2, hours: 10)),
        status: 'Pending',
        empIdx: 3,
        svcIdx: 3,
      ),
      build(
        id: 5,
        when: today
            .subtract(const Duration(days: 1))
            .add(const Duration(hours: 12)),
        status: 'Approved',
        empIdx: 4,
        svcIdx: 1,
        process: 'Completed',
        notes: 'Dün 12:00 tamamlandı',
      ),
      build(
        id: 6,
        when: today
            .subtract(const Duration(days: 3))
            .add(const Duration(hours: 14)),
        status: 'Canceled',
        empIdx: 5,
        svcIdx: 7,
        process: 'Completed',
      ),
      Appointment(
        appointmentId: 7,
        randevuId: 'demo-7',
        calisanId: 10,
        customerName: 'Ayşe Yılmaz',
        employeeName: emps[0].fullName,
        serviceName: svcs[2].serviceName,
        process: 'Waiting',
        totalPrice: svcs[2].servicePrice,
        appointmentDateTime: today.add(const Duration(hours: 13)),
        approvalStatus: 'Pending',
        notes: 'HIIT seansı',
        customerPhone: '0555 987 65 43',
        customerEmail: 'ayse@orion.com',
        createdAt: now,
        updatedAt: now,
      ),
      Appointment(
        appointmentId: 8,
        randevuId: 'demo-8',
        calisanId: 13,
        customerName: 'Mehmet Demir',
        employeeName: emps[3].fullName,
        serviceName: svcs[3].serviceName,
        process: 'In Progress',
        totalPrice: svcs[3].servicePrice,
        appointmentDateTime: today.add(const Duration(days: 3, hours: 11)),
        approvalStatus: 'Approved',
        notes: 'CrossFit',
        customerPhone: '0555 444 55 66',
        customerEmail: 'mehmet@orion.com',
        createdAt: now,
        updatedAt: now,
      ),
    ];
  }

  static List<Appointment> appointmentsForCustomer(String email) {
    final target =
        email.trim().isEmpty ? demoEmail : email.trim().toLowerCase();
    return appointments(forEmail: target)
        .where((a) => (a.customerEmail ?? '').toLowerCase() == target)
        .toList();
  }

  static List<Appointment> profileAppointments(String email) {
    final mine = appointmentsForCustomer(email);
    if (mine.isNotEmpty) return mine;
    return appointmentsForCustomer(demoEmail);
  }

  static Map<String, dynamic> isletme() => {
        'isim': 'Orion Gym & Fitness',
        'aciklama': 'Spor Salonu & Fitness',
        'banner_url': image('orion-banner'),
        // Navbar logo ağ görseli kullanmaz — tema rengi fitness ikonu
        'logo_url': '',
        'arka_plan_url': image('orion-bg'),
        'telefon': '0212 555 00 00',
        'email': 'info@orion.com',
        'adres': 'Caddebostan, İstanbul',
        'hero_tagline': 'Güçlü Hisset, Orion Ol',
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
        image('event-1'),
        image('event-2'),
        image('event-3'),
      ];

  static List<Map<String, String>> team() => [
        {
          'name': 'Zeynep Yılmaz',
          'role': 'Kişisel Antrenör',
          'image': image('team-1'),
        },
        {
          'name': 'Can Öztürk',
          'role': 'Yoga Eğitmeni',
          'image': image('team-2'),
        },
        {
          'name': 'Elif Demir',
          'role': 'Pilates Uzmanı',
          'image': image('team-3'),
        },
        {
          'name': 'Burak Kaya',
          'role': 'CrossFit Coach',
          'image': image('team-4'),
        },
        {
          'name': 'Selin Arslan',
          'role': 'Grup Fitness',
          'image': image('team-5'),
        },
        {
          'name': 'Mert Çelik',
          'role': 'Boks Antrenörü',
          'image': image('team-6'),
        },
      ];

  static Map<String, String> icerikBlok() => {
        'why_1_title': 'Uzman Antrenörler',
        'why_1_desc':
            'Sertifikalı ekibimizle hedeflerine güvenli ve hızlı ulaş',
        'why_2_title': 'Modern Ekipman',
        'why_2_desc': 'Son teknoloji aletlerle verimli antrenman',
        'why_3_title': 'Esnek Program',
        'why_3_desc': 'Sabah-akşam dersleri, her seviyeye uygun paketler',
        'currency': 'TL',
        'dialog_close': 'Kapat',
        'hero_subtitle': 'Spor Salonu & Fitness',
        'hero_tagline': 'Güçlü Hisset, Orion Ol',
        'brand_name': 'Orion',
        'btn_login': 'Giriş Yap',
        'btn_learn_more': 'Daha Fazla',
      };

  static List<String> serviceCategories() => [
        'All',
        'Fitness',
        'Wellness',
        'Kardiyo',
        'Spor',
      ];

  /// Admin performans kartları
  static List<Map<String, dynamic>> employeePerformanceMaps() {
    final emps = employees();
    return List.generate(emps.length, (i) {
      final e = emps[i];
      final seed = (e.id ?? (i + 1));
      return {
        'employeeName': e.fullName,
        'dailyEarnings': 800.0 + seed * 120.0,
        'totalEarnings': 12000.0 + seed * 1500.0,
        'efficiency': 72.0 + (seed % 20),
        'date': DateTime.now(),
        'appointmentsCompleted': 3 + (seed % 8),
        'averageRating': 4.0 + (seed % 10) / 10.0,
        'pendingAppointments': seed % 3,
      };
    });
  }

  /// Admin ana sayfa — son aktiviteler
  static List<Map<String, dynamic>> recentActivities() {
    final appts = appointments();
    return appts.take(8).map((a) {
      final status = a.approvalStatus.toLowerCase();
      String title;
      if (status.contains('pending')) {
        title = 'Yeni randevu: ${a.customerName} — ${a.serviceName}';
      } else if (status.contains('cancel')) {
        title = 'İptal: ${a.customerName} — ${a.serviceName}';
      } else if (a.process.toLowerCase().contains('complete')) {
        title = 'Tamamlandı: ${a.customerName} — ${a.serviceName}';
      } else {
        title = 'Onaylandı: ${a.customerName} — ${a.serviceName}';
      }
      return {
        'title': title,
        'time': _relativeDemo(a.appointmentDateTime),
        'ts': a.appointmentDateTime,
        'status': a.approvalStatus,
      };
    }).toList();
  }

  static Map<String, int> todaySummary() {
    final now = DateTime.now();
    final today = DateTime(now.year, now.month, now.day);
    final tomorrow = today.add(const Duration(days: 1));
    final todays = appointments().where((a) {
      final d = a.appointmentDateTime;
      return !d.isBefore(today) && d.isBefore(tomorrow);
    }).toList();
    final pending = todays
        .where((a) => a.approvalStatus.toLowerCase().contains('pending'))
        .length;
    final approved = todays
        .where((a) => a.approvalStatus.toLowerCase().contains('approved'))
        .length;
    return {
      'total': todays.length,
      'pending': pending,
      'completed': approved,
    };
  }

  static String _relativeDemo(DateTime dt) {
    final diff = DateTime.now().difference(dt);
    if (diff.isNegative) {
      final ahead = dt.difference(DateTime.now());
      if (ahead.inHours < 1) return '${ahead.inMinutes} dk sonra';
      if (ahead.inHours < 24) return '${ahead.inHours} saat sonra';
      return '${ahead.inDays} gün sonra';
    }
    if (diff.inMinutes < 60) return '${diff.inMinutes} dk önce';
    if (diff.inHours < 24) return '${diff.inHours} saat önce';
    return '${diff.inDays} gün önce';
  }

  static String _title(String s) {
    if (s.isEmpty) return demoName;
    return s[0].toUpperCase() + s.substring(1);
  }
}
