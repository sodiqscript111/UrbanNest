import { useState } from 'react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Search } from 'lucide-react';

const Header = () => {
    const [location, setLocation] = useState('');
    const [checkIn, setCheckIn] = useState(null);
    const [checkOut, setCheckOut] = useState(null);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleFindNest = () => {
        navigate('/home');
    };

    const handleSearch = (e) => {
        e.preventDefault();
        if (!location || !checkIn || !checkOut) {
            setError('Please fill all fields');
            return;
        }
        setError('');
        const query = {
            location,
            checkIn: checkIn ? checkIn.toISOString().split('T')[0] : '',
            checkOut: checkOut ? checkOut.toISOString().split('T')[0] : '',
        };
        console.log('Search query:', query);
        // TODO: Navigate to search results with query
        navigate(`/home?location=${encodeURIComponent(location)}&checkIn=${query.checkIn}&checkOut=${query.checkOut}`);
    };

    // Animation variants
    const cardVariants = {
        hidden: { opacity: 0, y: 20 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: 'easeOut' } },
    };

    const imageVariants = {
        hidden: { opacity: 0 },
        visible: { opacity: 1, transition: { duration: 0.6, ease: 'easeOut' } },
    };

    const errorVariants = {
        hidden: { opacity: 0, y: -10 },
        visible: { opacity: 1, y: 0, transition: { duration: 0.3, ease: 'easeOut' } },
    };

    return (
        <header className="bg-white font-inter">
            <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
                <div className="flex flex-col gap-4 sm:gap-6">
                    {/* Image */}
                    <motion.div
                        className="relative h-48 sm:h-64 md:h-[32rem]"
                        variants={imageVariants}
                        initial="hidden"
                        animate="visible"
                    >
                        <img
                            src="https://i.imgur.com/jRpocK5.jpeg"
                            alt="Modern apartment in Lagos"
                            className="w-full h-full object-cover rounded-lg"
                            loading="lazy"
                        />
                    </motion.div>

                    {/* Search Form */}
                    <motion.div
                        className="bg-white shadow-lg rounded-2xl p-6 sm:p-8 max-w-full"
                        variants={cardVariants}
                        initial="hidden"
                        animate="visible"
                    >
                        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-4 sm:mb-6">
                            Find Your Perfect Nest
                        </h1>
                        <form onSubmit={handleSearch} className="space-y-4">
                            <div>
                                <label htmlFor="location" className="block text-base font-medium text-gray-900">
                                    Location
                                </label>
                                <input
                                    type="text"
                                    id="location"
                                    value={location}
                                    onChange={(e) => setLocation(e.target.value)}
                                    placeholder="e.g., Lagos, Nigeria"
                                    className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-base text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 transition-all duration-200"
                                    aria-describedby="location-description"
                                />
                                <p id="location-description" className="text-sm text-gray-500 mt-1">
                                    Enter a city or neighborhood in Lagos.
                                </p>
                            </div>
                            <div className="flex flex-col sm:flex-row sm:space-x-4 space-y-4 sm:space-y-0">
                                <div className="flex-1">
                                    <label htmlFor="check-in" className="block text-base font-medium text-gray-900">
                                        Check-In
                                    </label>
                                    <DatePicker
                                        id="check-in"
                                        selected={checkIn}
                                        onChange={(date) => setCheckIn(date)}
                                        minDate={new Date()}
                                        placeholderText="Select date"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-base text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 transition-all duration-200"
                                        aria-label="Check-in date"
                                        wrapperClassName="w-full"
                                        popperClassName="z-50"
                                    />
                                </div>
                                <div className="flex-1">
                                    <label htmlFor="check-out" className="block text-base font-medium text-gray-900">
                                        Check-Out
                                    </label>
                                    <DatePicker
                                        id="check-out"
                                        selected={checkOut}
                                        onChange={(date) => {
                                            if (checkIn && date < checkIn) {
                                                setError('Check-out must be after check-in');
                                                return;
                                            }
                                            setCheckOut(date);
                                            setError('');
                                        }}
                                        minDate={checkIn || new Date()}
                                        placeholderText="Select date"
                                        className="mt-1 w-full px-4 py-3 border border-gray-300 rounded-md text-base text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-gray-500 transition-all duration-200"
                                        aria-label="Check-out date"
                                        wrapperClassName="w-full"
                                        popperClassName="z-50"
                                    />
                                </div>
                            </div>
                            {error && (
                                <motion.p
                                    className="text-sm text-red-500 text-center"
                                    variants={errorVariants}
                                    initial="hidden"
                                    animate="visible"
                                >
                                    {error}
                                </motion.p>
                            )}
                            <motion.button
                                type="submit"
                                className="w-full bg-white border border-black text-black font-semibold py-3 px-4 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 flex items-center justify-center gap-2"
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                transition={{ duration: 0.2 }}
                                aria-label="Search for listings"
                            >
                                <Search className="h-5 w-5" />
                                <span className="hidden sm:inline">Search</span>
                            </motion.button>
                            <motion.button
                                type="button"
                                onClick={handleFindNest}
                                className="w-full bg-black text-white font-semibold py-3 px-4 rounded-md hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-gray-500"
                                whileHover={{ scale: 1.05 }}
                                whileTap={{ scale: 0.95 }}
                                transition={{ duration: 0.2 }}
                                aria-label="Find your nest"
                            >
                                Find Your Nest
                            </motion.button>
                        </form>
                    </motion.div>
                </div>
            </div>
        </header>
    );
};

export default Header;