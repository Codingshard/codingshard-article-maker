// Ensure highlight.js is configured with the languages you want to support.
// This should be done BEFORE Quill is initialized.
// The 'go' language is specifically added here.
// ...existing code...

// Make sure Go is available in the Quill code block language dropdown
if (window.Quill && window.hljs) {
  // Quill 2.x: set the languages for the syntax module
  Quill.register('modules/syntax', Quill.import('modules/syntax'), true);
  Quill.import('modules/syntax').languages = [
    { id: 'go', label: 'Go' },
    { id: 'javascript', label: 'JavaScript' },
    { id: 'python', label: 'Python' },
    { id: 'java', label: 'Java' },
    { id: 'c++', label: 'C++' },
    { id: 'csharp', label: 'C#' },
    { id: 'php', label: 'PHP' },
    { id: 'ruby', label: 'Ruby' },
    { id: 'bash', label: 'Bash' },
    { id: 'html', label: 'HTML' },
    { id: 'css', label: 'CSS' },
    { id: 'sql', label: 'SQL' },
    { id: 'xml', label: 'XML' },
    // Add more as needed
  ];
}

hljs.configure({
  // Add all the languages you expect to use in your editor.
  // 'go' is included here.
  languages: [
    'go',
    'javascript',
    'html',
    'css',
    'python',
    'java',
    'bash',
    'c++',
    'csharp',
    'diff',
    'markdown',
    'php',
    'ruby',
    'sql',
    'xml', // Often useful for HTML/XML
    // You can add more languages as needed from highlight.js/lib/languages/
  ]
});

// Initialize Quill editor
const quill = new Quill("#editor", {
  theme: "snow",
  modules: {
    syntax: {
      anguages: [
        { id: 'go', label: 'Go' },
        { id: 'javascript', label: 'JavaScript' },
        { id: 'python', label: 'Python' },
        { id: 'java', label: 'Java' },
        { id: 'c++', label: 'C++' },
        { id: 'csharp', label: 'C#' },
        { id: 'php', label: 'PHP' },
        { id: 'ruby', label: 'Ruby' },
        { id: 'bash', label: 'Bash' },
        { id: 'html', label: 'HTML' },
        { id: 'css', label: 'CSS' },
        { id: 'sql', label: 'SQL' },
        { id: 'xml', label: 'XML' },
        // Add more as needed
      ]
    },
    toolbar: [
      [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
      ['bold', 'italic', 'underline', 'strikethrough'],
      ['link', 'image', 'video'],
      ['blockquote', 'code-block'], // Ensure 'code-block' is in your toolbar
      [{ 'list': 'ordered'}, { 'list': 'bullet' }],
      ['clean']
    ]
  }
});

// 1. Get a reference to the publish button using its new ID
const publishButton = document.getElementById("publishButton");

// 2. Add an event listener for the 'click' event
publishButton.addEventListener("click", async () => {
  // 3. Get the HTML content from the Quill editor
  const htmlContent = quill.getSemanticHTML();
  const articleName = document.getElementById('articleName').value;

  // 4. Construct a JSON object for the server
  const payload = {
    htmlContent: htmlContent,
    articleName: articleName,
  };

  // 5. Send the POST request to your Go server
  try {
    const response = await fetch("/save-article", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    // 6. Handle the server's response
    const result = await response.json();
    if (response.ok) {
      console.log("Article published successfully!", result);
      // Use a custom message box instead of alert()
      showMessageBox("Article published successfully!");
    } else {
      console.error("Failed to publish article:", result.error);
      // Use a custom message box instead of alert()
      showMessageBox("Failed to publish article: " + result.error);
    }
  } catch (error) {
    console.error("Network error:", error);
    // Use a custom message box instead of alert()
    showMessageBox("A network error occurred.");
  }
});

// --- Custom Message Box Implementation (replaces alert()) ---
function showMessageBox(message) {
  const messageBoxId = 'custom-message-box';
  let messageBox = document.getElementById(messageBoxId);

  if (!messageBox) {
    messageBox = document.createElement('div');
    messageBox.id = messageBoxId;
    messageBox.style.cssText = `
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background-color: #333;
      color: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3);
      z-index: 1000;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 15px;
      font-family: 'Inter', sans-serif;
      max-width: 80%;
      text-align: center;
    `;
    document.body.appendChild(messageBox);
  }

  messageBox.innerHTML = `
    <p>${message}</p>
    <button id="message-box-ok-button" style="
      background-color: #007bff;
      color: white;
      border: none;
      padding: 10px 20px;
      border-radius: 5px;
      cursor: pointer;
      font-size: 1em;
      transition: background-color 0.2s ease;
    ">OK</button>
  `;

  messageBox.querySelector('#message-box-ok-button').onclick = () => {
    messageBox.remove();
  };

  messageBox.style.display = 'flex'; // Show the message box
}

