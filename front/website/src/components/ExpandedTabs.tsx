/** @jsxImportSource react */
import React, { useState, useRef, useEffect } from "react";
import { AnimatePresence, motion } from "framer-motion";

const ProductIcon = ({ className = "w-5 h-5" }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <path d="M6 9h12M6 9a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2M6 9v10a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V9M10 14h4M10 14v4M14 14v4" />
  </svg>
);

const HowIcon = ({ className = "w-5 h-5" }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <circle cx="12" cy="12" r="10" />
    <path d="M12 6v6l4 2" />
  </svg>
);

const PriceIcon = ({ className = "w-5 h-5" }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <line x1="12" y1="1" x2="12" y2="23" />
    <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
  </svg>
);

const IntegrationIcon = ({ className = "w-5 h-5" }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <path d="M12 5v14M5 12h14M7.5 9.5a2.5 2.5 0 1 1 5 0 2.5 2.5 0 0 1-5 0M11.5 14.5a2.5 2.5 0 1 1 5 0 2.5 2.5 0 0 1-5 0" />
  </svg>
);

const TrackingIcon = ({ className = "w-5 h-5" }) => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}>
    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
    <circle cx="12" cy="10" r="3" />
  </svg>
);

interface Tab {
  title: string;
  iconName: 'product' | 'how' | 'price' | 'integration' | 'tracking';
  href: string;
  type?: never;
}

interface Separator {
  type: "separator";
  title?: never;
  iconName?: never;
  href?: never;
}

type TabItem = Tab | Separator;

interface ExpandedTabsProps {
  tabs: TabItem[];
  className?: string;
  onChange?: (index: number | null) => void;
}

const spanVariants = {
  initial: { width: 0, opacity: 0 },
  animate: {
    width: "auto",
    opacity: 1,
    transition: { delay: 0.05, duration: 0.2, ease: "easeOut" as const },
  },
  exit: {
    width: 0,
    opacity: 0,
    transition: { duration: 0.1, ease: "easeIn" as const },
  },
};

function ExpandedTabs({ tabs, className, onChange }: ExpandedTabsProps) {
  const [selected, setSelected] = useState<number | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const currentPath = window.location.pathname;

    const activeIndex = tabs.findIndex((tab) => {
      if (tab.type === "separator") return false;
      const tabPath = tab.href.replace(/\/$/, '');
      const normalizedCurrent = currentPath.replace(/\/$/, '');
      return normalizedCurrent === tabPath || normalizedCurrent.endsWith(tabPath);
    });

    if (activeIndex !== -1) {
      setSelected(activeIndex);
    }
  }, [tabs]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        const currentPath = window.location.pathname;
        const activeIndex = tabs.findIndex((tab) => {
          if (tab.type === "separator") return false;
          const tabPath = tab.href.replace(/\/$/, '');
          const normalizedCurrent = currentPath.replace(/\/$/, '');
          return normalizedCurrent === tabPath || normalizedCurrent.endsWith(tabPath);
        });

        if (activeIndex === -1) {
          setSelected(null);
          if (onChange) onChange(null);
        }
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [onChange, tabs]);

  const handleSelect = (index: number) => {
    setSelected(index);
    if (onChange) onChange(index);
  };

  const SeparatorComponent = () => (
    <div className="h-7 w-px bg-purple-200" aria-hidden="true" />
  );

  return (
    <div
      ref={containerRef}
      className={`flex items-center gap-1 rounded-full border border-purple-300 bg-white/95 p-1 shadow-md backdrop-blur-sm ${className || ""}`}
    >
      {tabs.map((tab, index) => {
        if (tab.type === "separator") {
          return <SeparatorComponent key={`separator-${index}`} />;
        }

        const iconMap = {
          product: ProductIcon,
          how: HowIcon,
          price: PriceIcon,
          integration: IntegrationIcon,
          tracking: TrackingIcon,
        };

        const Icon = iconMap[tab.iconName];
        const isSelected = selected === index;

        return (
          <a
            key={tab.title}
            href={tab.href}
            onMouseEnter={() => handleSelect(index)}
            onMouseLeave={() => setSelected(null)}
            className={`relative z-10 flex items-center rounded-full px-4 py-2 text-sm font-medium transition-colors focus:outline-none cursor-pointer ${
              isSelected ? "text-white" : "text-gray-700 hover:text-[#8B5CF6]"
            }`}
          >
            {isSelected && (
              <motion.div
                layoutId="pill"
                className="absolute inset-0 z-0 rounded-full bg-[#8B5CF6] shadow-sm"
                transition={{ type: "spring", stiffness: 500, damping: 40 }}
              />
            )}

            <span className="relative z-10 flex items-center gap-2">
              <Icon className="h-5 w-5 flex-shrink-0" />
              <span className="overflow-hidden whitespace-nowrap">{tab.title}</span>
            </span>
          </a>
        );
      })}
    </div>
  );
}

export default ExpandedTabs;
