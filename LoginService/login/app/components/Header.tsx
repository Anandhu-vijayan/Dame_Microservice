"use client";

import { useState } from "react";
import Image from "next/image";
import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/solid";

export default function Header() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      {/* Header (Fixed to Top) */}
      <header className="fixed top-0 left-0 w-full z-50 bg-white shadow-md dark:bg-gray-900" style={{ background: '#ebe7cd', height: '100px' }}>
        <div className="container mx-auto flex items-center justify-between px-4 py-2">
          {/* Left Side: Banner Image (Responsive) */}
          <div className="flex-shrink-0">
            <Image
              src="/assets/banner2.png"
              alt="Banner"
              width={800}
              height={50}
              className="rounded-lg object-contain md:w-[600px] w-[350px] md:ml-0 ml-[-20px]" // Adjust width and shift left
            />
          </div>

          {/* Desktop Menu */}
          <nav className="hidden md:flex space-x-6">
            <a href="#" className="text-gray-700 dark:text-gray-900 hover:text-[#238636]">Home</a>
            <a href="#" className="text-gray-700 dark:text-gray-900 hover:text-[#238636]">How to Apply</a>
            <a href="#" className="text-gray-700 dark:text-gray-900 hover:text-[#238636]">Prospectus</a>
            <a href="#" className="text-gray-700 dark:text-gray-900 hover:text-[#238636]">Notifications</a>
            <a href="#" className="text-gray-700 dark:text-gray-900 hover:text-[#238636]">Contact</a>
          </nav>

          {/* Mobile Menu Button */}
          <button
            className="md:hidden text-gray-700 dark:text-gray-900 focus:outline-none"
            onClick={() => setIsOpen(!isOpen)}
          >
            {isOpen ? <XMarkIcon className="w-8 h-8" /> : <Bars3Icon className="w-8 h-8" />}
          </button>
        </div>

        {/* Mobile Menu Dropdown */}
        {isOpen && (
          <div className="md:hidden bg-white dark:white p-4 space-y-2">
            <a href="#" className="block text-gray-700 dark:text-gray-900 hover:text-blue-500">Home</a>
            <a href="#" className="block text-gray-700 dark:text-gray-900 hover:text-blue-500">Services</a>
            <a href="#" className="block text-gray-700 dark:text-gray-900 hover:text-blue-500">About</a>
            <a href="#" className="block text-gray-700 dark:text-gray-900 hover:text-blue-500">Contact</a>
            <a href="#" className="block text-gray-700 dark:text-gray-900 hover:text-blue-500">Login</a>
          </div>
        )}
      </header>

      {/* Adjust Homepage Content to Avoid Overlap */}
      <div className="pt-[90px]">
        {/* The rest of the page content goes here */}
      </div>
    </>
  );
}
