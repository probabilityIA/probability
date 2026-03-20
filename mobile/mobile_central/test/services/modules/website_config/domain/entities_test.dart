import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/website_config/domain/entities.dart';

void main() {
  group('WebsiteConfigData', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'template': 'modern',
        'show_hero': true,
        'show_about': true,
        'show_featured_products': true,
        'show_full_catalog': false,
        'show_testimonials': true,
        'show_location': false,
        'show_contact': true,
        'show_social_media': true,
        'show_whatsapp': false,
        'hero_content': {'title': 'Welcome', 'subtitle': 'Shop now'},
        'about_content': {'text': 'About us'},
        'testimonials_content': [
          {'name': 'John', 'text': 'Great!'},
          {'name': 'Jane', 'text': 'Excellent!'},
        ],
        'location_content': {'address': '123 Main St'},
        'contact_content': {'email': 'info@test.com'},
        'social_media_content': {'instagram': '@test'},
        'whatsapp_content': {'number': '+573001234567'},
      };

      final config = WebsiteConfigData.fromJson(json);

      expect(config.id, 1);
      expect(config.businessId, 5);
      expect(config.template, 'modern');
      expect(config.showHero, true);
      expect(config.showAbout, true);
      expect(config.showFeaturedProducts, true);
      expect(config.showFullCatalog, false);
      expect(config.showTestimonials, true);
      expect(config.showLocation, false);
      expect(config.showContact, true);
      expect(config.showSocialMedia, true);
      expect(config.showWhatsapp, false);
      expect(config.heroContent, isNotNull);
      expect(config.heroContent!['title'], 'Welcome');
      expect(config.aboutContent, isNotNull);
      expect(config.aboutContent!['text'], 'About us');
      expect(config.testimonialsContent, isNotNull);
      expect(config.testimonialsContent!.length, 2);
      expect(config.testimonialsContent![0]['name'], 'John');
      expect(config.locationContent, isNotNull);
      expect(config.contactContent, isNotNull);
      expect(config.socialMediaContent, isNotNull);
      expect(config.whatsappContent, isNotNull);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final config = WebsiteConfigData.fromJson(json);

      expect(config.id, 0);
      expect(config.businessId, 0);
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
      expect(config.heroContent, isNull);
      expect(config.aboutContent, isNull);
      expect(config.testimonialsContent, isNull);
      expect(config.locationContent, isNull);
      expect(config.contactContent, isNull);
      expect(config.socialMediaContent, isNull);
      expect(config.whatsappContent, isNull);
    });

    test('fromJson handles null content fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'template': 'basic',
        'show_hero': false,
        'show_about': false,
        'show_featured_products': false,
        'show_full_catalog': false,
        'show_testimonials': false,
        'show_location': false,
        'show_contact': false,
        'show_social_media': false,
        'show_whatsapp': false,
        'hero_content': null,
        'about_content': null,
        'testimonials_content': null,
        'location_content': null,
        'contact_content': null,
        'social_media_content': null,
        'whatsapp_content': null,
      };

      final config = WebsiteConfigData.fromJson(json);

      expect(config.heroContent, isNull);
      expect(config.aboutContent, isNull);
      expect(config.testimonialsContent, isNull);
      expect(config.locationContent, isNull);
      expect(config.contactContent, isNull);
      expect(config.socialMediaContent, isNull);
      expect(config.whatsappContent, isNull);
    });

    test('fromJson handles empty testimonials list', () {
      final json = {
        'testimonials_content': [],
      };

      final config = WebsiteConfigData.fromJson(json);

      expect(config.testimonialsContent, isNotNull);
      expect(config.testimonialsContent, isEmpty);
    });
  });

  group('UpdateWebsiteConfigDTO', () {
    test('toJson includes all non-null fields', () {
      final dto = UpdateWebsiteConfigDTO(
        template: 'modern',
        showHero: true,
        showAbout: false,
        showFeaturedProducts: true,
        showFullCatalog: false,
        showTestimonials: true,
        showLocation: false,
        showContact: true,
        showSocialMedia: true,
        showWhatsapp: false,
        heroContent: {'title': 'New Title'},
        aboutContent: {'text': 'New about'},
        testimonialsContent: [
          {'name': 'Test', 'text': 'Review'},
        ],
        locationContent: {'address': 'New addr'},
        contactContent: {'email': 'new@test.com'},
        socialMediaContent: {'ig': '@new'},
        whatsappContent: {'number': '+57'},
      );

      final json = dto.toJson();

      expect(json['template'], 'modern');
      expect(json['show_hero'], true);
      expect(json['show_about'], false);
      expect(json['show_featured_products'], true);
      expect(json['show_full_catalog'], false);
      expect(json['show_testimonials'], true);
      expect(json['show_location'], false);
      expect(json['show_contact'], true);
      expect(json['show_social_media'], true);
      expect(json['show_whatsapp'], false);
      expect(json['hero_content'], isNotNull);
      expect(json['hero_content']['title'], 'New Title');
      expect(json['about_content'], isNotNull);
      expect(json['testimonials_content'], isList);
      expect((json['testimonials_content'] as List).length, 1);
      expect(json['location_content'], isNotNull);
      expect(json['contact_content'], isNotNull);
      expect(json['social_media_content'], isNotNull);
      expect(json['whatsapp_content'], isNotNull);
    });

    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateWebsiteConfigDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateWebsiteConfigDTO(
        template: 'basic',
        showHero: true,
      );

      final json = dto.toJson();

      expect(json.length, 2);
      expect(json['template'], 'basic');
      expect(json['show_hero'], true);
      expect(json.containsKey('show_about'), false);
      expect(json.containsKey('hero_content'), false);
    });

    test('toJson excludes null content fields', () {
      final dto = UpdateWebsiteConfigDTO(
        showHero: true,
        heroContent: null,
      );

      final json = dto.toJson();

      expect(json.containsKey('show_hero'), true);
      expect(json.containsKey('hero_content'), false);
    });

    test('toJson handles empty testimonials list', () {
      final dto = UpdateWebsiteConfigDTO(
        testimonialsContent: [],
      );

      final json = dto.toJson();

      expect(json['testimonials_content'], isEmpty);
    });
  });
}
