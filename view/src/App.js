import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import SignUp from './components/SignUp';
import Login from './components/Login';
import Greeting from './components/Greeting';

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/sign-up" element={<SignUp />} />
                <Route path="/sign-up" element={<Login />} />
                <Route path="/auth/greeting" element={<Greeting />} />
            </Routes>
        </Router>
    );
}

export default App;
