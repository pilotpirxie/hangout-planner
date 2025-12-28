import { useEffect, useState } from "react";

export type ScreenSize = "mobile" | "tablet" | "desktop";

export const useResponsive = () => {
  const getScreenSize = (width: number): ScreenSize => {
    if (width <= 576) return "mobile";
    if (width <= 992) return "tablet";
    return "desktop";
  };

  const [screenSize, setScreenSize] = useState<ScreenSize>(
    getScreenSize(window.innerWidth)
  );

  useEffect(() => {
    const handleResize = () => {
      setScreenSize(getScreenSize(window.innerWidth));
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return {
    screenSize,
    isMobile: screenSize === "mobile",
    isTablet: screenSize === "tablet",
    isDesktop: screenSize === "desktop",
  };
};
