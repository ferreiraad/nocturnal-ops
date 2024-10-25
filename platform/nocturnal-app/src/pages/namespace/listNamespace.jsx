import React, { useEffect, useState } from 'react';
import M from 'materialize-css';

const ListNamespacePage = () => {
  const [namespaces, setNamespaces] = useState([]);
  const [newNamespace, setNewNamespace] = useState("");
  const [error, setError] = useState("");

  const fetchNamespaces = async () => {
    try {
      console.log('REFRESH');
      const response = await fetch('http://localhost:9000/namespaces');
      const data = await response.json();
      setNamespaces(data.namespaces);
    } catch (error) {
      console.error("Error fetching namespaces:", error);
    }
  };

  useEffect(() => {
    // Fetch namespaces on component mount
    fetchNamespaces();

    // Initialize Materialize modal
    const modalElement = document.querySelector('#addNamespaceModal');
    M.Modal.init(modalElement);
  }, []);

  // Function to handle Add button click
  const handleAddNamespace = async () => {
    if (!newNamespace.trim()) {
      setError("Namespace cannot be empty");
      return;
    }

    try {
      const response = await fetch('http://localhost:9000/namespaces', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ namespace: newNamespace }),
      });

      if (response.ok) {
        setNewNamespace("");
        setError("");
        fetchNamespaces(); // Refresh the list

        // Close the modal
        const modalInstance = M.Modal.getInstance(document.getElementById('addNamespaceModal'));
        modalInstance.close();
      } else {
        console.error("Error adding namespace:", response.statusText);
      }
    } catch (error) {
      console.error("Error adding namespace:", error);
    }
  };

  return (
    <div className="container">
      <h2>Namespaces</h2>

      {/* Buttons */}
      <button onClick={fetchNamespaces} className="btn">Refresh</button>
      <button data-target="addNamespaceModal" className="btn modal-trigger">Add</button>

      {/* Namespaces Table */}
      <table className="striped">
        <thead>
          <tr>
            <th>#</th>
            <th>Namespace</th>
          </tr>
        </thead>
        <tbody>
          {namespaces.length > 0 ? (
            namespaces.map((namespace, index) => (
              <tr key={index}>
                <td>{index + 1}</td>
                <td>{namespace}</td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="2">No namespaces available</td>
            </tr>
          )}
        </tbody>
      </table>

      {/* Add Namespace Modal */}
      <div id="addNamespaceModal" className="modal">
        <div className="modal-content">
          <h4>Add Namespace</h4>
          <input
            type="text"
            value={newNamespace}
            onChange={(e) => setNewNamespace(e.target.value)}
            placeholder="Enter namespace"
          />
          {error && <p style={{ color: "red" }}>{error}</p>}
        </div>
        <div className="modal-footer">
          <button onClick={handleAddNamespace} className="btn">Submit</button>
          <button className="modal-close btn red">Close</button>
        </div>
      </div>
    </div>
  );
};

export default ListNamespacePage;
