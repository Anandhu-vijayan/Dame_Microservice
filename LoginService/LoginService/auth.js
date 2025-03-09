import jwt from "jsonwebtoken";
import dotenv from "dotenv";

dotenv.config(); // Load environment variables

// Function to generate JWT token
export const generateToken = (userId) => {
    return jwt.sign({ userId }, process.env.JWT_SECRET, { expiresIn: process.env.JWT_EXPIRY });
};

// Function to verify JWT token
export const verifyToken = (token) => {
    return jwt.verify(token, process.env.JWT_SECRET);
};
