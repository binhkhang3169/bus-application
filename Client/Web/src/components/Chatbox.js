/** @format */

import React from "react";

// Placeholder function for when you implement the chat logic
const handleChatboxClick = () => {
  alert("Chức năng chat sẽ được phát triển sau!");
};

function Chatbox() {
  return (
    <button
      onClick={handleChatboxClick}
      className="fixed bottom-8 right-8 bg-blue-600 text-white w-16 h-16 rounded-full flex items-center justify-center shadow-lg hover:bg-blue-700 transition-colors duration-300 ease-in-out z-50"
      aria-label="Open chat"
    >
      {/* Icon SVG cho chatbox */}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className="h-8 w-8"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
        />
      </svg>
    </button>
  );
}

export default Chatbox;