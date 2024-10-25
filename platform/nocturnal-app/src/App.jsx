import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import HeaderComponent from './components/Header';
// import FooterComponent from './components/Footer';
import './App.scss';
import HomePage from './pages/home';
import AboutPage from './pages/about';
import ListNamespacePage from './pages/namespace/listNamespace';


const App = () => {
  return (
    <Router>
      <div className="App">
        <HeaderComponent />

        <main className='main-content'>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/about" element={<AboutPage />} />
            <Route path="/namespace/list" element={<ListNamespacePage />} />
          </Routes>
        </main>

        {/* <FooterComponent /> */}
      </div>
    </Router>
  );
};

export default App;
