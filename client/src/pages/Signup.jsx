import { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Lock, User, Phone } from 'lucide-react';

export default function Signup() {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        phone: '',
    });
    const [errors, setErrors] = useState({});
    const navigate = useNavigate();

    const BACKEND_URL = "https://urbannest-ybda.onrender.com";

    const emailRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z]{2,})+$/;
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
    const phoneRegex = /^\+?[1-9]\d{1,14}$/;

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));

        let error = '';
        if (name === 'email' && value && !emailRegex.test(value)) {
            error = 'Please enter a valid email address';
        } else if (name === 'password' && value && !passwordRegex.test(value)) {
            error = 'Password must be at least 8 characters, include uppercase, lowercase, number, and special character';
        } else if (name === 'firstName' && value && !/^[a-zA-Z]{1,50}$/.test(value)) {
            error = 'First name must be 1-50 letters';
        } else if (name === 'lastName' && value && !/^[a-zA-Z]{1,50}$/.test(value)) {
            error = 'Last name must be 1-50 letters';
        } else if (name === 'phone' && value && !phoneRegex.test(value)) {
            error = 'Please enter a valid phone number';
        }
        setErrors((prev) => ({ ...prev, [name]: error }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const newErrors = {};
        if (!emailRegex.test(formData.email)) {
            newErrors.email = 'Please enter a valid email address';
        }
        if (!passwordRegex.test(formData.password)) {
            newErrors.password = 'Password must be at least 8 characters, include uppercase, lowercase, number, and special character';
        }
        if (!/^[a-zA-Z]{1,50}$/.test(formData.firstName)) {
            newErrors.firstName = 'First name must be 1-50 letters';
        }
        if (!/^[a-zA-Z]{1,50}$/.test(formData.lastName)) {
            newErrors.lastName = 'Last name must be 1-50 letters';
        }
        if (!phoneRegex.test(formData.phone)) {
            newErrors.phone = 'Please enter a valid phone number';
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        const fullName = `${formData.firstName} ${formData.lastName}`;

        try {
            const res = await axios.post(`${BACKEND_URL}/signup`, {
                email: formData.email,
                password: formData.password,
                full_name: fullName,
                firstName: formData.firstName,
                lastName: formData.lastName,
                phone: formData.phone,
            });
            console.log('Signup response:', res.data);
            localStorage.setItem("token", res.data.token);
            alert('Signup successful! Please log in.');
            navigate('/login');
        } catch (err) {
            const errorMessage = err.response?.data?.error || 'Signup failed';
            console.error('Signup error:', err.response?.data || err);
            setErrors((prev) => ({ ...prev, submit: errorMessage }));
        }
    };

    return (
        <div className="min-h-screen bg-gray-100 flex items-center justify-center">
            <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
                <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                            Email
                        </label>
                        <div className="flex items-center border rounded-md">
                            <User className="w-5 h-5 text-gray-400 ml-2" />
                            <input
                                type="email"
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                className="w-full p-2 pl-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                placeholder="Enter your email"
                                required
                            />
                        </div>
                        {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email}</p>}
                    </div>
                    <div className="mb-4">
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                            Password
                        </label>
                        <div className="flex items-center border rounded-md">
                            <Lock className="w-5 h-5 text-gray-400 ml-2" />
                            <input
                                type="password"
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                className="w-full p-2 pl-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                placeholder="Enter your password"
                                required
                            />
                        </div>
                        {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password}</p>}
                    </div>
                    <div className="mb-4">
                        <label htmlFor="firstName" className="block text-sm font-medium text-gray-700">
                            First Name
                        </label>
                        <div className="flex items-center border rounded-md">
                            <User className="w-5 h-5 text-gray-400 ml-2" />
                            <input
                                type="text"
                                id="firstName"
                                name="firstName"
                                value={formData.firstName}
                                onChange={handleChange}
                                className="w-full p-2 pl-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                placeholder="Enter your first name"
                                required
                            />
                        </div>
                        {errors.firstName && <p className="text-red-500 text-sm mt-1">{errors.firstName}</p>}
                    </div>
                    <div className="mb-4">
                        <label htmlFor="lastName" className="block text-sm font-medium text-gray-700">
                            Last Name
                        </label>
                        <div className="flex items-center border rounded-md">
                            <User className="w-5 h-5 text-gray-400 ml-2" />
                            <input
                                type="text"
                                id="lastName"
                                name="lastName"
                                value={formData.lastName}
                                onChange={handleChange}
                                className="w-full p-2 pl-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                placeholder="Enter your last name"
                                required
                            />
                        </div>
                        {errors.lastName && <p className="text-red-500 text-sm mt-1">{errors.lastName}</p>}
                    </div>
                    <div className="mb-6">
                        <label htmlFor="phone" className="block text-sm font-medium text-gray-700">
                            Phone
                        </label>
                        <div className="flex items-center border rounded-md">
                            <Phone className="w-5 h-5 text-gray-400 ml-2" />
                            <input
                                type="tel"
                                id="phone"
                                name="phone"
                                value={formData.phone}
                                onChange={handleChange}
                                className="w-full p-2 pl-10 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                placeholder="Enter your phone number"
                                required
                            />
                        </div>
                        {errors.phone && <p className="text-red-500 text-sm mt-1">{errors.phone}</p>}
                    </div>
                    {errors.submit && <p className="text-red-500 text-sm mb-4">{errors.submit}</p>}
                    <button
                        type="submit"
                        className="w-full bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600 transition-colors"
                    >
                        Sign Up
                    </button>
                </form>
                <p className="mt-4 text-center">
                    Already have an account?{' '}
                    <a href="/login" className="text-blue-500 hover:underline">
                        Log in
                    </a>
                </p>
            </div>
        </div>
    );
}