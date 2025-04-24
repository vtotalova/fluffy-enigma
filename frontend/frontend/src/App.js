import React, { useState, useEffect } from 'react';
import './App.css';
import startupIllustration from './startup-illustration.png';
//Prompt: "Create a startup idea generator app that generates 
//random startup ideas. The app should have a simple and clean design.
//Use React for the frontend and CSS for styling.
//The app should also have a feature to delete ideas and generate new ones.";
// 🧪 Mock startup ideas
function App() {
  const [ideas, setIdeas] = useState([]);
  const [currentIdea, setCurrentIdea] = useState(null);
 

  // Fetch ideas from the backend https://www.geeksforgeeks.org/how-to-fetch-data-from-an-api-in-reactjs/
  useEffect(() => {
    fetch("http://localhost:8080/api/ideas")
      .then(res => res.json())
      .then(data => {
        console.log("Fetched ideas:", data); 
        setIdeas(data);
        setCurrentIdea(data[0] || null);
      })
      .catch(err => console.error("Error fetching ideas:", err));
  }, []);

  const generateIdea = () => {
    if (ideas.length <= 1 || !currentIdea) return;

    const currentIndex = ideas.findIndex(idea => idea.id === currentIdea.id);
    const nextIndex = (currentIndex + 1) % ideas.length;
    setCurrentIdea(ideas[nextIndex]);
  };

  const deleteIdea = async () => {
    if (!currentIdea) {
      console.log("No current idea to delete");
      return;
    }
  
    console.log("Attempting to delete:", currentIdea.id);
  
    try {
      const response = await fetch(`http://localhost:8080/api/ideas/${currentIdea.id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json'
        }
      });
  
      if (response.ok) {
        const updatedIdeas = ideas.filter(idea => idea.id !== currentIdea.id);
        setIdeas(updatedIdeas);
  
        if (updatedIdeas.length > 0) {
          const nextIdea = updatedIdeas[0];
          setCurrentIdea(nextIdea);
        } else {
          setCurrentIdea(null);
        }
      } else {
        console.error("Failed to delete idea. Status:", response.status);
      }
    } catch (err) {
      console.error("Delete error:", err);
    }
  };
  // //Prompt CoPilot: "Give me a function to delete an idea from the list of ideas based on the generateIdea()."
  // const deleteIdea = () => {
  //   if (!currentIdea) return;
  //   fetch(`http://localhost:8080/api/ideas/${currentIdea.id}`, {
  //     method: 'DELETE'
  //   })
  //     .then(res => {
  //       if (res.ok) {
  //         const updatedIdeas = ideas.filter(idea => idea.id !== currentIdea.id);
  //         setIdeas(updatedIdeas);
  //         setCurrentIdea(updatedIdeas[0] || null);
  //       } else {
  //         console.error("Failed to delete idea");
  //       }
  //     })
  //     .catch(err => console.error("Delete error:", err));
  // };

  const getIdeaText = () => {
    if (!currentIdea) return "No ideas available";
    return `${currentIdea.this} for ${currentIdea.that}`;
  };

  return (
    <div className="App">
      <div className="container">
        <div className="content-left"><center>
          <h1>Welcome to Fluffy!</h1>
          <h2>Startup Ideas You Didn't Know You Needed</h2>

          <div className="idea-text">
            {getIdeaText()}
          </div>
          
            <div className="button-group">
              <button className="generate-button" onClick={generateIdea}>
                 Generate Idea
              </button>
              {currentIdea && (
              <button className="delete-button" onClick={deleteIdea}>
              Delete Idea
              </button>
            )}
            </div>
        </center>
        </div>

        <div className="content-right">
          <div className="illustration">
            <img src={startupIllustration} alt="Startup illustration" />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
