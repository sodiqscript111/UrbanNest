import { NavLink } from 'react-router-dom';
import { useState } from 'react';
import { Menu, X } from 'lucide-react';

const Nav = () => {
    const [isOpen, setIsOpen] = useState(false);

    const toggleMenu = () => {
        setIsOpen(!isOpen);
    };

    return (
        <nav className="bg-black text-white">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    {/* Logo */}
                    <div className="flex-shrink-0">
                        <NavLink to="/" className="text-2xl font-bold">
                            UrbanNest
                        </NavLink>
                    </div>

                    {/* Desktop Links */}
                    <div className="hidden md:flex space-x-8">
                        <NavLink
                            to="/"
                            className={({ isActive }) =>
                                `text-base font-medium hover:underline ${isActive ? 'underline' : ''}`
                            }
                        >
                            Nest a Place
                        </NavLink>
                        <NavLink
                            to="/create-listing"
                            className={({ isActive }) =>
                                `text-base font-medium hover:underline ${isActive ? 'underline' : ''}`
                            }
                        >
                            List a Nest
                        </NavLink>
                        <NavLink
                            to="/signup"
                            className={({ isActive }) =>
                                `text-base font-medium hover:underline ${isActive ? 'underline' : ''}`
                            }
                        >
                            Sign Up
                        </NavLink>
                    </div>

                    {/* Mobile Menu Button */}
                    <div className="md:hidden">
                        <button
                            onClick={toggleMenu}
                            className="p-2 rounded-md text-gray-400 hover:text-white focus:outline-none"
                            aria-label="Toggle menu"
                        >
                            {isOpen ? <X size={24} /> : <Menu size={24} />}
                        </button>
                    </div>
                </div>

                {/* Mobile Menu */}
                {isOpen && (
                    <div className="md:hidden">
                        <div className="px-2 pt-2 pb-3 space-y-1">
                            <NavLink
                                to="/"
                                className={({ isActive }) =>
                                    `block px-3 py-2 text-base font-medium hover:bg-gray-800 ${isActive ? 'bg-gray-800' : ''}`
                                }
                                onClick={toggleMenu}
                            >
                                Nest a Place
                            </NavLink>
                            <NavLink
                                to="/create-listing"
                                className={({ isActive }) =>
                                    `block px-3 py-2 text-base font-medium hover:bg-gray-800 ${isActive ? 'bg-gray-800' : ''}`
                                }
                                onClick={toggleMenu}
                            >
                                List a Nest
                            </NavLink>
                            <NavLink
                                to="/signup"
                                className={({ isActive }) =>
                                    `block px-3 py-2 text-base font-medium hover:bg-gray-800 ${isActive ? 'bg-gray-800' : ''}`
                                }
                                onClick={toggleMenu}
                            >
                                Sign Up
                            </NavLink>
                        </div>
                    </div>
                )}
            </div>
        </nav>
    );
};

export default Nav