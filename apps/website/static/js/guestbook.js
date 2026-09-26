document.addEventListener("DOMContentLoaded", () => {
  const opinionsContainer = document.getElementById("opinions");

  fetch("http://localhost:8080/api/v1/guestbook")
    .then((response) => {
      // Check if the response is successful
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((opinions) => {
      opinions.forEach((opinion) => {
        const opinionDiv = document.createElement("div");
        opinionDiv.classList.add("single_opinion");

        // Create elements for the sanitized content
        const userNameHeading = document.createElement("h3");
        const userCommentParagraph = document.createElement("p");

        // Use textContent to prevent HTML injection
        userNameHeading.textContent = `${opinion.UserName} says...`;
        userCommentParagraph.textContent = opinion.UserComment;

        // Append the safe elements to the div
        opinionDiv.appendChild(userNameHeading);
        opinionDiv.appendChild(userCommentParagraph);

        opinionsContainer.appendChild(opinionDiv);
      });
    })
    .catch((error) =>
      console.error("Error fetching or processing opinions:", error)
    );
});

// Use an onload event on the window to ensure the DOM is fully loaded before trying to access the input element.
function update_post_url() {
  // Get the current URL from the browser's location object
  const currentUrl = window.location.href;

  // Log the URL to the browser's console for debugging purposes
  console.log("Current URL:", currentUrl);

  // Find the input element by its ID
  const urlInput = document.getElementById("post_url");

  // Set the value of the input field to the current URL
  if (urlInput) {
    urlInput.value = currentUrl;
  }
}
