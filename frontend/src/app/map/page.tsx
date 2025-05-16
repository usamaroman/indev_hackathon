"use client"

import { useEffect, useState } from 'react';
import { login } from '../auth';

import FirstFloorMap from "@/components/FirstFloorMap";

const floors = [1, 2, 3, 4, 5];

export default function HotelMapPage() {
 const [selectedFloor, setSelectedFloor] = useState(1);

const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const performLogin = async () => {
      try {
        const authData = await login();
        console.log('Login successful', authData);
        // Store token if needed
        localStorage.setItem('authToken', authData.access_token);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Login failed');
      } finally {
        setIsLoading(false);
      }
    };

    performLogin();
  }, []);


  const renderMap = () => {
    switch (selectedFloor) {
      case 1:
        return <FirstFloorMap />;
      case 2:
        return <div>2 этаж....</div>;
      default:
        return <div>Карта недоступна</div>;
    }
  };

if (isLoading) {
    return <div>Logging in...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  } 



  return (
    <main className="p-6">
      <header className="p-6">
        <h1 className="text-2xl font-bold">Карта Отеля</h1>
      </header>

      <div className="flex justify-center space-x-2 mb-4">
        {floors.map((floor) => (
          <button
            key={floor}
            onClick={() => setSelectedFloor(floor)}
            className={`px-4 py-2 border rounded ${
              selectedFloor === floor
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-800 hover:bg-gray-200'
            }`}
          >
            {floor}
          </button>
        ))}
      </div>

      <div className="flex justify-center">
        {renderMap()}
      </div>
    </main>
  );
}
