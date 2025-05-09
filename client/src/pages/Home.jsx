import PropertyCard from "../component/Propertycard.jsx";
import { useEffect, useState } from "react";

export default function Home() {
  const [listings, setListings] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

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
      const url = '/listings';
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
        const filteredListings = (data.listings || [])
            .sort((a, b) => b.price - a.price)
            .slice(0, 4);
        setListings(filteredListings);
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
    return <div className="text-gray-600 p-6 max-w-7xl mx-auto">Loading listings...</div>;
  }

  if (error) {
    return <div className="text-red-500 p-6 max-w-7xl mx-auto bg-white rounded-md shadow-md">{error}</div>;
  }

  return (
      <div className="bg-gray-100 py-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight leading-tight mb-8">
            Discover Your Perfect Stay
          </h1>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
            {listings.map((listing) => (
                <div
                    key={listing.id}
                    className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300"
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
                </div>
            ))}
          </div>
          {listings.length === 0 && !error && (
              <div className="mt-8 text-center text-gray-600">
                <p>No listings available at the moment. Please check back later.</p>
              </div>
          )}
        </div>
      </div>
  );
}