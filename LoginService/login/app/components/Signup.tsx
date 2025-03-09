"use client";

import { useState, useRef } from "react";
import Link from "next/link";
import { FiUpload } from "react-icons/fi";

export default function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [userImage, setUserImage] = useState<File | null>(null);
  const [aadharCards, setAadharCards] = useState<File[]>([]); // Multiple Aadhar PDFs
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  const fileInputRef1 = useRef<HTMLInputElement | null>(null);
  const fileInputRef2 = useRef<HTMLInputElement | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setAadharCards(Array.from(event.target.files));
    }
  };
  

  const validateForm = () => {
    let errors: { [key: string]: string } = {};

    if (!/^[A-Za-z .]+$/.test(name)) {
      errors.name = "Name can only contain letters, spaces, and dots.";
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "Enter a valid email address.";
    }

    if (!/^\d{10}$/.test(phone)) {
      errors.phone = "Phone number must be exactly 10 digits.";
    }

    let passwordErrors = [];

    if (!/.{8,}/.test(password)) {
      passwordErrors.push("Must be at least 8 characters.");
    }
    if (!/[a-z]/.test(password)) {
      passwordErrors.push("Must include a lowercase letter.");
    }
    if (!/[A-Z]/.test(password)) {
      passwordErrors.push("Must include an uppercase letter.");
    }
    if (!/\d/.test(password)) {
      passwordErrors.push("Must include a number.");
    }
    if (!/[@$!%*?&]/.test(password)) {
      passwordErrors.push("Must include a special character (@$!%*?&).");
    }

    if (passwordErrors.length > 0) {
      errors.password = passwordErrors.join(" ");
    }

    setErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
  
    if (!validateForm()) return;
  
    const formData = new FormData();
    formData.append("name", name);
    formData.append("email", email);
    formData.append("phone", phone);
    formData.append("password", password);
    if (userImage) formData.append("userImage", userImage);
    
    // Append multiple PDFs correctly
    aadharCards.forEach((file) => {
      formData.append("pdfFiles", file); // Using the same key multiple times
    });
  
    try {
      console.log(formData);
      const response = await fetch("http://localhost:8080/register", {
        method: "POST",
        body: formData,
      });
  
      if (response.ok) {
        alert("User registered successfully!");
        setName("");
        setEmail("");
        setPhone("");
        setPassword("");
        setUserImage(null);
        setAadharCards([]);
        setErrors({});
        
        // Reset file input by changing key
        fileInputRef1.current?.value && (fileInputRef1.current.value = "");
        fileInputRef2.current?.value && (fileInputRef2.current.value = "");
      } else if (response.status === 409) {
        alert("User already registered");
        setName("");
        setEmail("");
        setPhone("");
        setPassword("");
        setUserImage(null);
        setAadharCards([]);
        setErrors({});
        
        fileInputRef1.current?.value && (fileInputRef1.current.value = "");
        fileInputRef2.current?.value && (fileInputRef2.current.value = "");
      } else {
        console.log(response);
        alert("Failed to register user.");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("An error occurred. Please try again.");
    }
  };
  

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-white px-4 sm:px-6 lg:px-8 pt-20">
      <div className="w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg mt-8">
        <h2 className="text-2xl font-bold text-center text-gray-900">
          Candidate Registration
        </h2>
        <form onSubmit={handleSubmit} className="mt-4" encType="multipart/form-data">
          {/* Name Field */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
            {errors.name && <p className="text-red-500 text-sm">{errors.name}</p>}
          </div>

          {/* Email Field */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Email ID</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
            {errors.email && <p className="text-red-500 text-sm">{errors.email}</p>}
          </div>

          {/* Phone Number Field */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Phone Number</label>
            <input
              type="tel"
              maxLength={10}
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
            {errors.phone && <p className="text-red-500 text-sm">{errors.phone}</p>}
          </div>

          {/* Password Field */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
            {errors.password && <p className="text-red-500 text-sm">{errors.password}</p>}
          </div>

          {/* Upload User Image */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Upload Photo</label>
            < div  className = "flex items-center border rounded-lg bg-gray-100 px-4 py-2" >
      <FiUpload className="text-gray-500 mr-2" />
            <input type="file" accept="image/*" ref={fileInputRef1} onChange={(e) => setUserImage(e.target.files?.[0] || null)} className="w-full bg-white text-gray-900 focus:outline-none" required />
          </div>
</div>
          {/* Upload Aadhar Card */}
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">Upload Aadhar Card (Multiple PDFs)</label>
            < div  className = "flex items-center border rounded-lg bg-gray-100 px-4 py-2" >
      <FiUpload className="text-gray-500 mr-2" />
            <input type="file" accept="application/pdf" ref={fileInputRef2} multiple onChange={handleFileChange} className="w-full bg-white text-gray-900 focus:outline-none" required />
          </div>
</div>
          <button type="submit" className="w-full bg-green-600 text-white py-2 rounded-lg hover:bg-green-700 transition-all">
            Register
          </button>
        </form>
        < p className = "mt-4 text-lg sm:text-base text-gray-600 text-center" >
        Already have an account ? { " "}
          < Link  href = "/" className = "text-blue-500 hover:underline font-semibold" > Login Here </Link>
            </p>
      </div>
    </div>
  );
}
