"use client";

import { useState } from "react";
import Link from "next/link";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await fetch("http://localhost:5000/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();
      if (response.ok) {
        console.log("Login Successful:", data);
        // Redirect user or store token in localStorage
      } else {
        setError(data.message || "Invalid credentials");
      }
    } catch (err) {
      console.error("Login error:", err);
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-white px-4 sm:px-6 lg:px-8 pt-20 sm:pt-22">
      <h2 className="text-xl sm:text-2xl md:text-3xl font-bold text-gray-900 text-center leading-snug">
        ONLINE APPLICATION PORTAL FOR <br className="hidden sm:block" />
        GOVERNMENT PARAMEDICAL COURSES
      </h2>

      <div className="w-full max-w-md bg-white p-6 sm:p-8 rounded-2xl shadow-lg mt-8 sm:mt-10">
        <h2 className="text-2xl font-bold text-center text-gray-900">
          Candidate Login
        </h2>
        <form onSubmit={handleSubmit} className="mt-4">
          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">
              Username
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
          </div>

          <div className="mb-3">
            <label className="block text-gray-700 font-medium mb-1">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg text-gray-900 bg-gray-100 focus:outline-none focus:ring-2 focus:ring-green-500"
              required
            />
          </div>

          {error && <p className="text-red-600 text-center">{error}</p>}

          <div className="flex justify-between items-center mb-4">
            <a href="#" className="text-sm text-blue-500 hover:underline">
              Forgot Password?
            </a>
          </div>

          <button
            type="submit"
            className="w-full bg-green-600 text-white py-2 rounded-lg hover:bg-green-700 transition-all"
            disabled={loading}
          >
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <p className="mt-4 text-lg sm:text-base text-gray-600 text-center">
          New User?{" "}
          <Link href="/signup" className="text-blue-500 hover:underline font-semibold">
            Register Now
          </Link>
        </p>
      </div>
    </div>
  );
}
