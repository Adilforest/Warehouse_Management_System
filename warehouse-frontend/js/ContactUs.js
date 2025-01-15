document.addEventListener('DOMContentLoaded', function () {
    const contactForm = document.getElementById('contactForm');
    const responseMessage = document.getElementById('responseMessage');

    contactForm.addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent the default form submission

        // Create a FormData object to collect form data
        const formData = new FormData(contactForm);

        // Send the data to the backend
        fetch('http://localhost:8080/api/contact', {
            method: 'POST',
            body: formData,
        })
            .then(response => {
                if (response.status === 429) {
                    // If the request limit is exceeded
                    throw new Error('Too Many Requests');
                }
                if (!response.ok) {
                    // If the response is not OK, throw an error
                    throw new Error('Network response was not ok');
                }
                return response.json(); // Parse JSON from the response
            })
            .then(data => {
                if (data.success) {
                    // Successful submission
                    responseMessage.textContent = 'Message sent successfully!';
                    responseMessage.classList.remove('error');
                    responseMessage.classList.add('success');
                    contactForm.reset(); // Clear the form
                } else {
                    // Server-side error
                    responseMessage.textContent = 'Error sending message. Please try again.';
                    responseMessage.classList.remove('success');
                    responseMessage.classList.add('error');
                }
            })
            .catch(error => {
                // Handle errors
                console.error('Error:', error);

                if (error.message === 'Too Many Requests') {
                    // If the request limit is exceeded
                    responseMessage.textContent = 'You have sent too many requests. Please try again in 15 seconds.';
                } else {
                    // Other errors (e.g., network issues)
                    responseMessage.textContent = 'An error occurred while sending. Please check your internet connection.';
                }

                responseMessage.classList.remove('success');
                responseMessage.classList.add('error');
            });
    });
});