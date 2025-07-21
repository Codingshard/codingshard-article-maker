// Initialize Quill editor
const quill = new Quill("#editor", {
  theme: "snow",
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
      alert("Article published successfully!");
    } else {
      console.error("Failed to publish article:", result.error);
      alert("Failed to publish article: " + result.error);
    }
  } catch (error) {
    console.error("Network error:", error);
    alert("A network error occurred.");
  }
});
