class PublicBusiness {
  final int id;
  final String name;
  final String code;
  final String description;
  final String logoUrl;
  final String primaryColor;
  final String secondaryColor;
  final String tertiaryColor;
  final String quaternaryColor;
  final String navbarImageUrl;
  final WebsiteConfig? websiteConfig;
  final List<PublicProduct> featuredProducts;

  PublicBusiness({
    required this.id,
    required this.name,
    required this.code,
    required this.description,
    required this.logoUrl,
    required this.primaryColor,
    required this.secondaryColor,
    required this.tertiaryColor,
    required this.quaternaryColor,
    required this.navbarImageUrl,
    this.websiteConfig,
    required this.featuredProducts,
  });

  factory PublicBusiness.fromJson(Map<String, dynamic> json) {
    return PublicBusiness(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      description: json['description'] ?? '',
      logoUrl: json['logo_url'] ?? '',
      primaryColor: json['primary_color'] ?? '',
      secondaryColor: json['secondary_color'] ?? '',
      tertiaryColor: json['tertiary_color'] ?? '',
      quaternaryColor: json['quaternary_color'] ?? '',
      navbarImageUrl: json['navbar_image_url'] ?? '',
      websiteConfig: json['website_config'] != null
          ? WebsiteConfig.fromJson(json['website_config'])
          : null,
      featuredProducts: (json['featured_products'] as List<dynamic>?)
              ?.map((e) => PublicProduct.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class WebsiteConfig {
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
  final HeroContent? heroContent;
  final AboutContent? aboutContent;
  final List<Testimonial>? testimonialsContent;
  final LocationContent? locationContent;
  final ContactContent? contactContent;
  final SocialMediaContent? socialMediaContent;
  final WhatsAppContent? whatsappContent;

  WebsiteConfig({
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

  factory WebsiteConfig.fromJson(Map<String, dynamic> json) {
    return WebsiteConfig(
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
          ? HeroContent.fromJson(json['hero_content'])
          : null,
      aboutContent: json['about_content'] != null
          ? AboutContent.fromJson(json['about_content'])
          : null,
      testimonialsContent: (json['testimonials_content'] as List<dynamic>?)
          ?.map((e) => Testimonial.fromJson(e))
          .toList(),
      locationContent: json['location_content'] != null
          ? LocationContent.fromJson(json['location_content'])
          : null,
      contactContent: json['contact_content'] != null
          ? ContactContent.fromJson(json['contact_content'])
          : null,
      socialMediaContent: json['social_media_content'] != null
          ? SocialMediaContent.fromJson(json['social_media_content'])
          : null,
      whatsappContent: json['whatsapp_content'] != null
          ? WhatsAppContent.fromJson(json['whatsapp_content'])
          : null,
    );
  }
}

class HeroContent {
  final String? title;
  final String? subtitle;
  final String? ctaText;
  final String? backgroundImage;

  HeroContent({
    this.title,
    this.subtitle,
    this.ctaText,
    this.backgroundImage,
  });

  factory HeroContent.fromJson(Map<String, dynamic> json) {
    return HeroContent(
      title: json['title'],
      subtitle: json['subtitle'],
      ctaText: json['cta_text'],
      backgroundImage: json['background_image'],
    );
  }
}

class AboutContent {
  final String? text;
  final String? image;
  final String? mission;
  final String? vision;

  AboutContent({
    this.text,
    this.image,
    this.mission,
    this.vision,
  });

  factory AboutContent.fromJson(Map<String, dynamic> json) {
    return AboutContent(
      text: json['text'],
      image: json['image'],
      mission: json['mission'],
      vision: json['vision'],
    );
  }
}

class Testimonial {
  final String name;
  final String text;
  final int? rating;
  final String? avatar;

  Testimonial({
    required this.name,
    required this.text,
    this.rating,
    this.avatar,
  });

  factory Testimonial.fromJson(Map<String, dynamic> json) {
    return Testimonial(
      name: json['name'] ?? '',
      text: json['text'] ?? '',
      rating: json['rating'],
      avatar: json['avatar'],
    );
  }
}

class LocationContent {
  final double? lat;
  final double? lng;
  final String? address;
  final String? hours;

  LocationContent({
    this.lat,
    this.lng,
    this.address,
    this.hours,
  });

  factory LocationContent.fromJson(Map<String, dynamic> json) {
    return LocationContent(
      lat: json['lat']?.toDouble(),
      lng: json['lng']?.toDouble(),
      address: json['address'],
      hours: json['hours'],
    );
  }
}

class ContactContent {
  final String? email;
  final String? phone;
  final bool? formEnabled;
  final List<ContactPerson>? contacts;

  ContactContent({
    this.email,
    this.phone,
    this.formEnabled,
    this.contacts,
  });

  factory ContactContent.fromJson(Map<String, dynamic> json) {
    return ContactContent(
      email: json['email'],
      phone: json['phone'],
      formEnabled: json['form_enabled'],
      contacts: (json['contacts'] as List<dynamic>?)
          ?.map((e) => ContactPerson.fromJson(e))
          .toList(),
    );
  }
}

class ContactPerson {
  final String name;
  final String role;
  final String phone;

  ContactPerson({
    required this.name,
    required this.role,
    required this.phone,
  });

  factory ContactPerson.fromJson(Map<String, dynamic> json) {
    return ContactPerson(
      name: json['name'] ?? '',
      role: json['role'] ?? '',
      phone: json['phone'] ?? '',
    );
  }
}

class SocialMediaContent {
  final String? facebook;
  final String? instagram;
  final String? twitter;
  final String? tiktok;

  SocialMediaContent({
    this.facebook,
    this.instagram,
    this.twitter,
    this.tiktok,
  });

  factory SocialMediaContent.fromJson(Map<String, dynamic> json) {
    return SocialMediaContent(
      facebook: json['facebook'],
      instagram: json['instagram'],
      twitter: json['twitter'],
      tiktok: json['tiktok'],
    );
  }
}

class WhatsAppContent {
  final String? number;
  final String? message;
  final bool? showFloatingButton;

  WhatsAppContent({
    this.number,
    this.message,
    this.showFloatingButton,
  });

  factory WhatsAppContent.fromJson(Map<String, dynamic> json) {
    return WhatsAppContent(
      number: json['number'],
      message: json['message'],
      showFloatingButton: json['show_floating_button'],
    );
  }
}

class PublicProduct {
  final String id;
  final String name;
  final String description;
  final String shortDescription;
  final double price;
  final double? compareAtPrice;
  final String currency;
  final String imageUrl;
  final List<String>? images;
  final String sku;
  final int stockQuantity;
  final String category;
  final String brand;
  final bool isFeatured;
  final String createdAt;

  PublicProduct({
    required this.id,
    required this.name,
    required this.description,
    required this.shortDescription,
    required this.price,
    this.compareAtPrice,
    required this.currency,
    required this.imageUrl,
    this.images,
    required this.sku,
    required this.stockQuantity,
    required this.category,
    required this.brand,
    required this.isFeatured,
    required this.createdAt,
  });

  factory PublicProduct.fromJson(Map<String, dynamic> json) {
    return PublicProduct(
      id: json['id']?.toString() ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      shortDescription: json['short_description'] ?? '',
      price: (json['price'] ?? 0).toDouble(),
      compareAtPrice: json['compare_at_price']?.toDouble(),
      currency: json['currency'] ?? 'COP',
      imageUrl: json['image_url'] ?? '',
      images: (json['images'] as List<dynamic>?)?.map((e) => e.toString()).toList(),
      sku: json['sku'] ?? '',
      stockQuantity: json['stock_quantity'] ?? 0,
      category: json['category'] ?? '',
      brand: json['brand'] ?? '',
      isFeatured: json['is_featured'] ?? false,
      createdAt: json['created_at'] ?? '',
    );
  }
}

class ContactFormDTO {
  final String name;
  final String? email;
  final String? phone;
  final String message;

  ContactFormDTO({
    required this.name,
    this.email,
    this.phone,
    required this.message,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'message': message,
    };
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    return json;
  }
}

class GetPublicCatalogParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final String? category;

  GetPublicCatalogParams({
    this.page,
    this.pageSize,
    this.search,
    this.category,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (category != null) params['category'] = category;
    return params;
  }
}
