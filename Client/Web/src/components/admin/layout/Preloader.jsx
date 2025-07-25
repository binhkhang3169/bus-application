// src/components/admin/layout/Preloader.jsx
import React from "react";

const Preloader = () => {
  return (
    <div
      x-show="loaded"
      x-init="window.addEventListener('DOMContentLoaded', () => {setTimeout(() => loaded = false, 500)})"
      class="fixed left-0 top-0 z-999999 flex h-screen w-screen items-center justify-center bg-white dark:bg-black"
    >
      <div class="h-16 w-16 animate-spin rounded-full border-4 border-solid border-brand-500 border-t-transparent"></div>
    </div>
  );
};

export default Preloader;
