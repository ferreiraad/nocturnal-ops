import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import M from 'materialize-css';
import './Header.scss';
import batImage from '../assets/images/nice-bat.png';

const HeaderComponent = () => {
  useEffect(() => {
    const sidenav = document.querySelectorAll('.sidenav');
    M.Sidenav.init(sidenav, { edge: 'right' });
  }, []);

  const openSidenav = () => {
    const sidenavInstance = M.Sidenav.getInstance(document.getElementById('mobile-demo'));
    sidenavInstance.open();
  };

  return (
    <>
      <nav>
        <div className="nav-wrapper">
          <Link to="/" className="brand-logo">NocturnalOps</Link>

          <ul className="right">
            <li><Link to="/">Home</Link></li>
            <li><Link to="/about">About</Link></li>
            <li>
              <Link to="#" onClick={openSidenav}>
                <i className="material-icons">menu</i>
              </Link>
            </li>
          </ul>
        </div>
      </nav>

      {/* Sidenav for mobile view with titled groups */}
      <ul className="sidenav" id="mobile-demo">
        {/* <li className="sidenav-title">Namespace</li> */}
        <li><Link to="/namespace/list">Namespace</Link></li>
        {/* <li><Link to="/namespace/create">Create</Link></li> */}

        {/* <li className="sidenav-title">kind</li>
        <li><Link to="/kind/list">List</Link></li>
        <li><Link to="/kind/delete">Delete</Link></li>


        <li className="sidenav-title">Entity</li>
        <li><Link to="/kind/list">List</Link></li>
        <li><Link to="/kind/get">Get</Link></li>
        <li><Link to="/kind/filter">Filter</Link></li>
        <li><Link to="/kind/create">Create</Link></li>
        <li><Link to="/kind/update">Update</Link></li>
        <li><Link to="/kind/delete">Delete</Link></li> */}
      </ul>
    </>
  );
};

export default HeaderComponent;
