import React, { useEffect, useState } from 'react';
import axios from 'axios';

const Greeting = () => {
    const [message, setMessage] = useState('');

    useEffect(() => {
        const fetchGreeting = async () => {
            const token = localStorage.getItem('jwt');
            try {
                const response = await axios.get('http://localhost:8080/auth/greeting', {
                    headers: {
                        Authorization: token,
                    },
                });
                setMessage(response.data.message);
            } catch (error) {
                console.error('Error fetching greeting:', error);
            }
        };

        fetchGreeting();
    }, []);

    return <div>{message}</div>;
};

export default Greeting;
