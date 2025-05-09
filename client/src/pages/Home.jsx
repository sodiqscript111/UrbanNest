import PropertyCard from "../component/Propertycard.jsx";
import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

const SkeletonLoader = () => {
  return (
      <motion.div
          className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3 }}
      >
        <div className="h-6 bg-gray-200 rounded w-1/2 mb-6 animate-pulse"></div>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6">
          {[...Array(4)].map((_, index) => (
              <motion.div
                  key={index}
                  className="bg-white rounded-lg shadow-md overflow-hidden animate-pulse"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3, delay: index * 0.1 }}
              >
                <div className="w-full h-40 sm:h-48 bg-gray-200"></div>
                <div className="p-3 sm:p-4">
                  <div className="h-5 bg-gray-200 rounded w-3/4 mb-2"></div>
                  <div className="h-4 bg-gray-200 rounded w-full mb-2"></div>
                  <div className="h-4 bg-gray-200 rounded w-5/6 mb-2"></div>
                  <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
                  <div className="h-5 bg-gray-200 rounded w-1/3"></div>
                </div>
              </motion.div>
          ))}
        </div>
      </motion.div>
  );
};

export default function Home() {
  const [listings, setListings] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  // Scroll to top on mount
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  useEffect(() => {
    async function fetchWithRetry(url, options, retries = 3, delay = 1000) {
      for (let i = 0; i < retries; i++) {
        try {
          console.log(`Attempt ${i + 1} - Fetching listings from:`, url);
          const res = await fetch(url, options);
          console.log('Response status:', res.status);
          if (!res.ok) {
            const text = await res.text();
            console.log('Response body:', text);
            throw new Error(`HTTP error! status: ${res.status}`);
          }
          const data = await res.json();
          console.log('Response data:', data);
          return data;
        } catch (err) {
          console.error(`Attempt ${i + 1} failed:`, err);
          if (i < retries - 1) {
            console.log(`Retrying after ${delay}ms...`);
            await new Promise((resolve) => setTimeout(resolve, delay));
          } else {
            throw err;
          }
        }
      }
    }

    async function fetchListings() {
      const url = 'https://urbannest-ybda.onrender.com/listings';
      try {
        const data = await fetchWithRetry(
            url,
            {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              },
            },
            3,
            1000
        );
        const sortedListings = (data.listings || []).sort((a, b) => b.price - a.price);
        setListings(sortedListings);
      } catch (err) {
        console.error('Failed to fetch listings:', err);
        setError(`Failed to load listings: ${err.message}`);
      } finally {
        setIsLoading(false);
      }
    }

    fetchListings();
  }, []);

  if (isLoading) {
    return <SkeletonLoader />;
  }

  if (error) {
    return (
        <motion.div
            className="text-red-500 p-4 sm:p-6 max-w-7xl mx-auto bg-white rounded-md shadow-md font-inter"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        >
          {error}
        </motion.div>
    );
  }

  return (
      <div className="bg-gray-100 py-8 sm:py-10 font-inter">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-6 sm:mb-8">
            <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 tracking-tight leading-tight">
              Discover Your Perfect Stay
            </h1>
          </div>
          <AnimatePresence>
            <motion.div
                className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:gap-6"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.3 }}
            >
              {listings.map((listing, index) => (
                  <motion.div
                      key={listing.id}
                      className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200"
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.3, delay: index * 0.1 }}
                      whileHover={{ scale: 1.03 }}
                      whileTap={{ scale: 0.95 }}
                  >
                    <PropertyCard
                        id={listing.id}
                        media={
                          listing.media && listing.media.length > 0
                              ? listing.media.map(
                                  (item) => `https://urbannestbucket.s3.eu-north-1.amazonaws.com/${item.media_url}`
                              )
                              : ["https://via.placeholder.com/600x400"]
                        }
                        name={listing.title}
                        rating="4.5 ★ (100)"
                        description={listing.description.substring(0, 100) + "..."}
                        availability={listing.is_available ? "Available" : "Booked"}
                        price={listing.price}
                        className="rounded-lg"
                    />
                  </motion.div>
              ))}
            </motion.div>
          </AnimatePresence>
          {listings.length === 0 && !error && (
              <motion.div
                  className="mt-6 sm:mt-8 text-center text-gray-600 text-base sm:text-lg font-inter"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ duration: 0.3 }}
              >
                <p>No listings available at the moment. Please check back later.</p>
              </motion.div>
          )}
        </div>
      </div>
  );
}