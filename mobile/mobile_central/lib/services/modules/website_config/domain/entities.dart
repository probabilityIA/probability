class WebsiteConfigData {
  final int id;
  final int businessId;
  final String template;
  final bool showHero;
  final bool showAbout;
  final bool showFeaturedProducts;
  final bool showFullCatalog;
  final bool showTestimonials;
  final bool showLocation;
  final bool showContact;
  final bool showSocialMedia;
  final bool showWhatsapp;
  final Map<String, dynamic>? heroContent;
  final Map<String, dynamic>? aboutContent;
  final List<Map<String, dynamic>>? testimonialsContent;
  final Map<String, dynamic>? locationContent;
  final Map<String, dynamic>? contactContent;
  final Map<String, dynamic>? socialMediaContent;
  final Map<String, dynamic>? whatsappContent;

  WebsiteConfigData({
    required this.id,
    required this.businessId,
    required this.template,
    required this.showHero,
    required this.showAbout,
    required this.showFeaturedProducts,
    required this.showFullCatalog,
    required this.showTestimonials,
    required this.showLocation,
    required this.showContact,
    required this.showSocialMedia,
    required this.showWhatsapp,
    this.heroContent,
    this.aboutContent,
    this.testimonialsContent,
    this.locationContent,
    this.contactContent,
    this.socialMediaContent,
    this.whatsappContent,
  });

  factory WebsiteConfigData.fromJson(Map<String, dynamic> json) {
    return WebsiteConfigData(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      template: json['template'] ?? '',
      showHero: json['show_hero'] ?? false,
      showAbout: json['show_about'] ?? false,
      showFeaturedProducts: json['show_featured_products'] ?? false,
      showFullCatalog: json['show_full_catalog'] ?? false,
      showTestimonials: json['show_testimonials'] ?? false,
      showLocation: json['show_location'] ?? false,
      showContact: json['show_contact'] ?? false,
      showSocialMedia: json['show_social_media'] ?? false,
      showWhatsapp: json['show_whatsapp'] ?? false,
      heroContent: json['hero_content'] != null
          ? Map<String, dynamic>.from(json['hero_content'])
          : null,
      aboutContent: json['about_content'] != null
          ? Map<String, dynamic>.from(json['about_content'])
          : null,
      testimonialsContent: (json['testimonials_content'] as List<dynamic>?)
          ?.map((e) => Map<String, dynamic>.from(e))
          .toList(),
      locationContent: json['location_content'] != null
          ? Map<String, dynamic>.from(json['location_content'])
          : null,
      contactContent: json['contact_content'] != null
          ? Map<String, dynamic>.from(json['contact_content'])
          : null,
      socialMediaContent: json['social_media_content'] != null
          ? Map<String, dynamic>.from(json['social_media_content'])
          : null,
      whatsappContent: json['whatsapp_content'] != null
          ? Map<String, dynamic>.from(json['whatsapp_content'])
          : null,
    );
  }
}

class UpdateWebsiteConfigDTO {
  final String? template;
  final bool? showHero;
  final bool? showAbout;
  final bool? showFeaturedProducts;
  final bool? showFullCatalog;
  final bool? showTestimonials;
  final bool? showLocation;
  final bool? showContact;
  final bool? showSocialMedia;
  final bool? showWhatsapp;
  final Map<String, dynamic>? heroContent;
  final Map<String, dynamic>? aboutContent;
  final List<Map<String, dynamic>>? testimonialsContent;
  final Map<String, dynamic>? locationContent;
  final Map<String, dynamic>? contactContent;
  final Map<String, dynamic>? socialMediaContent;
  final Map<String, dynamic>? whatsappContent;

  UpdateWebsiteConfigDTO({
    this.template,
    this.showHero,
    this.showAbout,
    this.showFeaturedProducts,
    this.showFullCatalog,
    this.showTestimonials,
    this.showLocation,
    this.showContact,
    this.showSocialMedia,
    this.showWhatsapp,
    this.heroContent,
    this.aboutContent,
    this.testimonialsContent,
    this.locationContent,
    this.contactContent,
    this.socialMediaContent,
    this.whatsappContent,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (template != null) json['template'] = template;
    if (showHero != null) json['show_hero'] = showHero;
    if (showAbout != null) json['show_about'] = showAbout;
    if (showFeaturedProducts != null) json['show_featured_products'] = showFeaturedProducts;
    if (showFullCatalog != null) json['show_full_catalog'] = showFullCatalog;
    if (showTestimonials != null) json['show_testimonials'] = showTestimonials;
    if (showLocation != null) json['show_location'] = showLocation;
    if (showContact != null) json['show_contact'] = showContact;
    if (showSocialMedia != null) json['show_social_media'] = showSocialMedia;
    if (showWhatsapp != null) json['show_whatsapp'] = showWhatsapp;
    if (heroContent != null) json['hero_content'] = heroContent;
    if (aboutContent != null) json['about_content'] = aboutContent;
    if (testimonialsContent != null) json['testimonials_content'] = testimonialsContent;
    if (locationContent != null) json['location_content'] = locationContent;
    if (contactContent != null) json['contact_content'] = contactContent;
    if (socialMediaContent != null) json['social_media_content'] = socialMediaContent;
    if (whatsappContent != null) json['whatsapp_content'] = whatsappContent;
    return json;
  }
}
