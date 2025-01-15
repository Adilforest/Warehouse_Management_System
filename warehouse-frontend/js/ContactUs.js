document.addEventListener('DOMContentLoaded', function () {
    const contactForm = document.getElementById('contactForm');
    const responseMessage = document.getElementById('responseMessage');

    contactForm.addEventListener('submit', function (event) {
        event.preventDefault(); // Предотвращаем стандартную отправку формы

        // Создаем объект FormData для сбора данных формы
        const formData = new FormData(contactForm);

        // Отправляем данные на бэкенд
        fetch('http://localhost:8080/api/contact', {
            method: 'POST',
            body: formData,
        })
            .then(response => {
                if (!response.ok) {
                    // Если ответ не OK, выбрасываем ошибку
                    throw new Error('Network response was not ok');
                }
                return response.json(); // Парсим JSON из ответа
            })
            .then(data => {
                if (data.success) {
                    // Успешная отправка
                    responseMessage.textContent = 'Сообщение успешно отправлено!';
                    responseMessage.classList.remove('error');
                    responseMessage.classList.add('success');
                    contactForm.reset(); // Очищаем форму
                } else {
                    // Ошибка на стороне сервера
                    responseMessage.textContent = 'Ошибка при отправке сообщения. Пожалуйста, попробуйте еще раз.';
                    responseMessage.classList.remove('success');
                    responseMessage.classList.add('error');
                }
            })
            .catch(error => {
                // Ошибка сети или другая ошибка
                console.error('Ошибка:', error);
                responseMessage.textContent = 'Произошла ошибка при отправке. Пожалуйста, проверьте подключение к интернету.';
                responseMessage.classList.remove('success');
                responseMessage.classList.add('error');
            });
    });
});