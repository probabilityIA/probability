import { TemplateComponents } from '../types';
import { PublicNav } from '../../components/PublicNav';
import { PublicFooter } from '../../components/PublicFooter';
import { WhatsAppButton } from '../../components/WhatsAppButton';
import { HeroSection } from '../../components/HeroSection';
import { AboutSection } from '../../components/AboutSection';
import { FeaturedProducts } from '../../components/FeaturedProducts';
import { TestimonialsSection } from '../../components/TestimonialsSection';
import { LocationSection } from '../../components/LocationSection';
import { ContactSection } from '../../components/ContactSection';
import { SocialMediaLinks } from '../../components/SocialMediaLinks';
import { PublicProductCard } from '../../components/PublicProductCard';

export const defaultTemplate: TemplateComponents = {
    Nav: PublicNav,
    Footer: PublicFooter,
    WhatsAppButton,
    HeroSection,
    AboutSection,
    FeaturedProducts,
    TestimonialsSection,
    LocationSection,
    ContactSection,
    SocialMediaLinks,
    ProductCard: PublicProductCard,
};
