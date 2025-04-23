import React, { useState } from 'react';
import './App.css';
import startupIllustration from './startup-illustration.png';
//Prompt: "Create a startup idea generator app that generates 
//random startup ideas. The app should have a simple and clean design.
//Use React for the frontend and CSS for styling.
//The app should also have a feature to delete ideas and generate new ones.";
// 🧪 Mock startup ideas
const mockIdeas = [
  { id: 1, this: "An AI assistant", that: "for cooking" },
  { id: 2, this: "A platform", that: "for dog yoga" },
  { id: 3, this: "A subscription box", that: "for plant lovers" }
];

function App() {
  const [ideas, setIdeas] = useState(mockIdeas);
  const [currentIdea, setCurrentIdea] = useState(mockIdeas[0]);

  const generateIdea = () => {
    if (ideas.length <= 1) return;

    const currentIndex = ideas.findIndex(idea => idea.id === currentIdea.id);
    const nextIndex = (currentIndex + 1) % ideas.length;
    setCurrentIdea(ideas[nextIndex]);
  };

  const deleteIdea = (id) => {
    const updatedIdeas = ideas.filter(idea => idea.id !== id);
    setIdeas(updatedIdeas);
    setCurrentIdea(updatedIdeas[0] || null);
  };

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
              <button className="delete-button" onClick={() => deleteIdea(currentIdea.id)}>
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
