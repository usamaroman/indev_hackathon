'use client'

import axios from 'axios';
import { useEffect, useState } from 'react';

interface Reservation {
  id: number;
  user_id: number;
  room_id: string;
  check_in: string;
  check_out: string;
  status: string;
  created_at: string;
  login: string;
}

export default function ReservationsPage() {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchReservations();
  }, []);

const fetchReservations = async () => {
    try {
      const response = await axios.get(
        'http://localhost:8080/v1/hotel/rooms/reservations/confirmed',
        {
          headers: {
            'accept': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
          }
        }
      );
      setReservations(response.data);
    } catch (err) {
      setError('Failed to fetch reservations');
      console.error('Error fetching reservations:', err);
    } finally {
      setLoading(false);
    }
  };

const handleCheckIn = async (reservationId: number) => {
    try {
      await axios.patch(
        `http://localhost:8080/v1/hotel/rooms/reservations/${reservationId}`,
        { status: 'checked_in' },
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
          }
        }
      );
      fetchReservations();
    } catch (err) {
      console.error('Error updating reservation status:', err);
      alert('Failed to update reservation status');
    }
  };

 const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Confirmed Reservations</h1>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {reservations.map((reservation) => (
          <li key={reservation.id} style={{ 
            border: '1px solid #ddd', 
            padding: '15px', 
            marginBottom: '10px',
            borderRadius: '5px'
          }}>
            <p><strong>Room:</strong> {reservation.room_id}</p>
            <p><strong>Guest:</strong> {reservation.login}</p>
            <p><strong>Check-in:</strong> {formatDate(reservation.check_in)}</p>
            <p><strong>Check-out:</strong> {formatDate(reservation.check_out)}</p>
            <p><strong>Status:</strong> {reservation.status}</p>
            {reservation.status === 'confirmed' && (
              <button 
                onClick={() => handleCheckIn(reservation.id)}
                style={{
                  backgroundColor: '#4CAF50',
                  color: 'white',
                  padding: '8px 16px',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                CHECK IN
              </button>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
