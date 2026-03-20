import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/publicsite/domain/entities.dart';

void main() {
  group('PublicBusiness', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'name': 'My Store',
        'code': 'mystore',
        'description': 'A great store',
        'logo_url': 'https://img.example.com/logo.png',
        'primary_color': '#FF0000',
        'secondary_color': '#00FF00',
        'tertiary_color': '#0000FF',
        'quaternary_color': '#FFFF00',
        'navbar_image_url': 'https://img.example.com/navbar.png',
        'website_config': {
          'template': 'modern',
          'show_hero': true,
          'show_about': true,
          'show_featured_products': true,
          'show_full_catalog': false,
          'show_testimonials': true,
          'show_location': false,
          'show_contact': true,
          'show_social_media': true,
          'show_whatsapp': true,
        },
        'featured_products': [
          {
            'id': '1',
            'name': 'Product 1',
            'description': 'Desc',
            'short_description': 'Short',
            'price': 29.99,
            'currency': 'COP',
            'image_url': 'https://img.example.com/1.jpg',
            'sku': 'SKU-1',
            'stock_quantity': 10,
            'category': 'Electronics',
            'brand': 'Brand',
            'is_featured': true,
            'created_at': '2026-01-01',
          },
        ],
      };

      final biz = PublicBusiness.fromJson(json);

      expect(biz.id, 1);
      expect(biz.name, 'My Store');
      expect(biz.code, 'mystore');
      expect(biz.description, 'A great store');
      expect(biz.logoUrl, 'https://img.example.com/logo.png');
      expect(biz.primaryColor, '#FF0000');
      expect(biz.secondaryColor, '#00FF00');
      expect(biz.tertiaryColor, '#0000FF');
      expect(biz.quaternaryColor, '#FFFF00');
      expect(biz.navbarImageUrl, 'https://img.example.com/navbar.png');
      expect(biz.websiteConfig, isNotNull);
      expect(biz.websiteConfig!.template, 'modern');
      expect(biz.featuredProducts.length, 1);
      expect(biz.featuredProducts[0].name, 'Product 1');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final biz = PublicBusiness.fromJson(json);

      expect(biz.id, 0);
      expect(biz.name, '');
      expect(biz.code, '');
      expect(biz.description, '');
      expect(biz.logoUrl, '');
      expect(biz.primaryColor, '');
      expect(biz.websiteConfig, isNull);
      expect(biz.featuredProducts, isEmpty);
    });
  });

  group('WebsiteConfig', () {
    test('fromJson parses all boolean flags correctly', () {
      final json = {
        'template': 'classic',
        'show_hero': true,
        'show_about': true,
        'show_featured_products': true,
        'show_full_catalog': true,
        'show_testimonials': true,
        'show_location': true,
        'show_contact': true,
        'show_social_media': true,
        'show_whatsapp': true,
      };

      final config = WebsiteConfig.fromJson(json);

      expect(config.template, 'classic');
      expect(config.showHero, true);
      expect(config.showAbout, true);
      expect(config.showFeaturedProducts, true);
      expect(config.showFullCatalog, true);
      expect(config.showTestimonials, true);
      expect(config.showLocation, true);
      expect(config.showContact, true);
      expect(config.showSocialMedia, true);
      expect(config.showWhatsapp, true);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final config = WebsiteConfig.fromJson(json);

      expect(config.template, '');
      expect(config.showHero, false);
      expect(config.showAbout, false);
      expect(config.showFeaturedProducts, false);
      expect(config.showFullCatalog, false);
      expect(config.showTestimonials, false);
      expect(config.showLocation, false);
      expect(config.showContact, false);
      expect(config.showSocialMedia, false);
      expect(config.showWhatsapp, false);
    });

    test('fromJson parses nested content objects', () {
      final json = {
        'template': 't',
        'show_hero': false,
        'show_about': false,
        'show_featured_products': false,
        'show_full_catalog': false,
        'show_testimonials': false,
        'show_location': false,
        'show_contact': false,
        'show_social_media': false,
        'show_whatsapp': false,
        'hero_content': {'title': 'Welcome', 'subtitle': 'Sub'},
        'about_content': {'text': 'About us', 'mission': 'Mission'},
        'testimonials_content': [
          {'name': 'John', 'text': 'Great!', 'rating': 5},
        ],
        'location_content': {'lat': 4.6, 'lng': -74.0, 'address': '123 St'},
        'contact_content': {'email': 'a@b.com', 'phone': '123'},
        'social_media_content': {'facebook': 'fb.com/store'},
        'whatsapp_content': {'number': '+57123456789'},
      };

      final config = WebsiteConfig.fromJson(json);

      expect(config.heroContent, isNotNull);
      expect(config.heroContent!.title, 'Welcome');
      expect(config.aboutContent, isNotNull);
      expect(config.aboutContent!.text, 'About us');
      expect(config.testimonialsContent, isNotNull);
      expect(config.testimonialsContent!.length, 1);
      expect(config.locationContent, isNotNull);
      expect(config.locationContent!.lat, 4.6);
      expect(config.contactContent, isNotNull);
      expect(config.contactContent!.email, 'a@b.com');
      expect(config.socialMediaContent, isNotNull);
      expect(config.socialMediaContent!.facebook, 'fb.com/store');
      expect(config.whatsappContent, isNotNull);
      expect(config.whatsappContent!.number, '+57123456789');
    });

    test('fromJson handles null content objects', () {
      final json = {
        'template': 't',
        'show_hero': false,
        'show_about': false,
        'show_featured_products': false,
        'show_full_catalog': false,
        'show_testimonials': false,
        'show_location': false,
        'show_contact': false,
        'show_social_media': false,
        'show_whatsapp': false,
      };

      final config = WebsiteConfig.fromJson(json);

      expect(config.heroContent, isNull);
      expect(config.aboutContent, isNull);
      expect(config.testimonialsContent, isNull);
      expect(config.locationContent, isNull);
      expect(config.contactContent, isNull);
      expect(config.socialMediaContent, isNull);
      expect(config.whatsappContent, isNull);
    });
  });

  group('HeroContent', () {
    test('fromJson parses all fields', () {
      final json = {
        'title': 'Welcome',
        'subtitle': 'To our store',
        'cta_text': 'Shop Now',
        'background_image': 'https://img.example.com/bg.jpg',
      };

      final hero = HeroContent.fromJson(json);

      expect(hero.title, 'Welcome');
      expect(hero.subtitle, 'To our store');
      expect(hero.ctaText, 'Shop Now');
      expect(hero.backgroundImage, 'https://img.example.com/bg.jpg');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final hero = HeroContent.fromJson(json);

      expect(hero.title, isNull);
      expect(hero.subtitle, isNull);
      expect(hero.ctaText, isNull);
      expect(hero.backgroundImage, isNull);
    });
  });

  group('AboutContent', () {
    test('fromJson parses all fields', () {
      final json = {
        'text': 'About us',
        'image': 'https://img.example.com/about.jpg',
        'mission': 'Our mission',
        'vision': 'Our vision',
      };

      final about = AboutContent.fromJson(json);

      expect(about.text, 'About us');
      expect(about.image, 'https://img.example.com/about.jpg');
      expect(about.mission, 'Our mission');
      expect(about.vision, 'Our vision');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final about = AboutContent.fromJson(json);

      expect(about.text, isNull);
      expect(about.image, isNull);
      expect(about.mission, isNull);
      expect(about.vision, isNull);
    });
  });

  group('Testimonial', () {
    test('fromJson parses all fields', () {
      final json = {
        'name': 'John',
        'text': 'Great product!',
        'rating': 5,
        'avatar': 'https://img.example.com/john.jpg',
      };

      final t = Testimonial.fromJson(json);

      expect(t.name, 'John');
      expect(t.text, 'Great product!');
      expect(t.rating, 5);
      expect(t.avatar, 'https://img.example.com/john.jpg');
    });

    test('fromJson uses defaults for required and handles null optional', () {
      final json = <String, dynamic>{};
      final t = Testimonial.fromJson(json);

      expect(t.name, '');
      expect(t.text, '');
      expect(t.rating, isNull);
      expect(t.avatar, isNull);
    });
  });

  group('LocationContent', () {
    test('fromJson parses all fields', () {
      final json = {
        'lat': 4.624335,
        'lng': -74.063644,
        'address': 'Calle 100 #10-30',
        'hours': 'Mon-Fri 9-5',
      };

      final loc = LocationContent.fromJson(json);

      expect(loc.lat, 4.624335);
      expect(loc.lng, -74.063644);
      expect(loc.address, 'Calle 100 #10-30');
      expect(loc.hours, 'Mon-Fri 9-5');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final loc = LocationContent.fromJson(json);

      expect(loc.lat, isNull);
      expect(loc.lng, isNull);
      expect(loc.address, isNull);
      expect(loc.hours, isNull);
    });
  });

  group('ContactContent', () {
    test('fromJson parses all fields including contacts list', () {
      final json = {
        'email': 'contact@store.com',
        'phone': '+57123456789',
        'form_enabled': true,
        'contacts': [
          {'name': 'Support', 'role': 'Manager', 'phone': '+57111'},
        ],
      };

      final contact = ContactContent.fromJson(json);

      expect(contact.email, 'contact@store.com');
      expect(contact.phone, '+57123456789');
      expect(contact.formEnabled, true);
      expect(contact.contacts, isNotNull);
      expect(contact.contacts!.length, 1);
      expect(contact.contacts![0].name, 'Support');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final contact = ContactContent.fromJson(json);

      expect(contact.email, isNull);
      expect(contact.phone, isNull);
      expect(contact.formEnabled, isNull);
      expect(contact.contacts, isNull);
    });
  });

  group('ContactPerson', () {
    test('fromJson parses all fields', () {
      final json = {'name': 'John', 'role': 'CEO', 'phone': '+57111'};
      final cp = ContactPerson.fromJson(json);

      expect(cp.name, 'John');
      expect(cp.role, 'CEO');
      expect(cp.phone, '+57111');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};
      final cp = ContactPerson.fromJson(json);

      expect(cp.name, '');
      expect(cp.role, '');
      expect(cp.phone, '');
    });
  });

  group('SocialMediaContent', () {
    test('fromJson parses all fields', () {
      final json = {
        'facebook': 'fb.com/store',
        'instagram': 'ig.com/store',
        'twitter': 'x.com/store',
        'tiktok': 'tiktok.com/@store',
      };

      final sm = SocialMediaContent.fromJson(json);

      expect(sm.facebook, 'fb.com/store');
      expect(sm.instagram, 'ig.com/store');
      expect(sm.twitter, 'x.com/store');
      expect(sm.tiktok, 'tiktok.com/@store');
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final sm = SocialMediaContent.fromJson(json);

      expect(sm.facebook, isNull);
      expect(sm.instagram, isNull);
      expect(sm.twitter, isNull);
      expect(sm.tiktok, isNull);
    });
  });

  group('WhatsAppContent', () {
    test('fromJson parses all fields', () {
      final json = {
        'number': '+57123456789',
        'message': 'Hello!',
        'show_floating_button': true,
      };

      final wa = WhatsAppContent.fromJson(json);

      expect(wa.number, '+57123456789');
      expect(wa.message, 'Hello!');
      expect(wa.showFloatingButton, true);
    });

    test('fromJson handles null fields', () {
      final json = <String, dynamic>{};
      final wa = WhatsAppContent.fromJson(json);

      expect(wa.number, isNull);
      expect(wa.message, isNull);
      expect(wa.showFloatingButton, isNull);
    });
  });

  group('PublicProduct', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'name': 'Widget',
        'description': 'A widget',
        'short_description': 'Short',
        'price': 29.99,
        'compare_at_price': 39.99,
        'currency': 'USD',
        'image_url': 'https://img.example.com/1.jpg',
        'images': ['https://img.example.com/1.jpg'],
        'sku': 'SKU-001',
        'stock_quantity': 100,
        'category': 'Electronics',
        'brand': 'BrandX',
        'is_featured': true,
        'created_at': '2026-01-01',
      };

      final p = PublicProduct.fromJson(json);

      expect(p.id, '42');
      expect(p.name, 'Widget');
      expect(p.description, 'A widget');
      expect(p.shortDescription, 'Short');
      expect(p.price, 29.99);
      expect(p.compareAtPrice, 39.99);
      expect(p.currency, 'USD');
      expect(p.imageUrl, 'https://img.example.com/1.jpg');
      expect(p.images, hasLength(1));
      expect(p.sku, 'SKU-001');
      expect(p.stockQuantity, 100);
      expect(p.category, 'Electronics');
      expect(p.brand, 'BrandX');
      expect(p.isFeatured, true);
      expect(p.createdAt, '2026-01-01');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final p = PublicProduct.fromJson(json);

      expect(p.id, '');
      expect(p.name, '');
      expect(p.description, '');
      expect(p.shortDescription, '');
      expect(p.price, 0.0);
      expect(p.currency, 'COP');
      expect(p.imageUrl, '');
      expect(p.sku, '');
      expect(p.stockQuantity, 0);
      expect(p.category, '');
      expect(p.brand, '');
      expect(p.isFeatured, false);
      expect(p.createdAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': '1',
        'name': 'P',
        'price': 1,
      };

      final p = PublicProduct.fromJson(json);

      expect(p.compareAtPrice, isNull);
      expect(p.images, isNull);
    });
  });

  group('ContactFormDTO', () {
    test('toJson includes required fields', () {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');

      final json = dto.toJson();

      expect(json['name'], 'John');
      expect(json['message'], 'Hello');
    });

    test('toJson includes optional fields when present', () {
      final dto = ContactFormDTO(
        name: 'John',
        email: 'john@example.com',
        phone: '+57123',
        message: 'Hello',
      );

      final json = dto.toJson();

      expect(json['email'], 'john@example.com');
      expect(json['phone'], '+57123');
    });

    test('toJson excludes null optional fields', () {
      final dto = ContactFormDTO(name: 'John', message: 'Hello');

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json.containsKey('email'), false);
      expect(json.containsKey('phone'), false);
    });
  });

  group('GetPublicCatalogParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetPublicCatalogParams(
        page: 1,
        pageSize: 20,
        search: 'widget',
        category: 'electronics',
      );

      final qp = params.toQueryParams();

      expect(qp['page'], 1);
      expect(qp['page_size'], 20);
      expect(qp['search'], 'widget');
      expect(qp['category'], 'electronics');
    });

    test('toQueryParams excludes null fields', () {
      final params = GetPublicCatalogParams(page: 1);

      final qp = params.toQueryParams();

      expect(qp.length, 1);
      expect(qp.containsKey('page'), true);
      expect(qp.containsKey('search'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetPublicCatalogParams();

      final qp = params.toQueryParams();

      expect(qp, isEmpty);
    });
  });
}
