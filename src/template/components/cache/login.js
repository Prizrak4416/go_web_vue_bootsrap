export default {
    data() {
        return {
          password: '',
          message: ""
        };
      },
      methods: {
        login() {
          fetch('/log', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ password: this.password }),
          })
          .then(response => {
              return response.json();
          })
          .then(data => {
            this.message = data.response
            // отправляем пароль на redirect
            localStorage.setItem('myData', JSON.stringify({ pass: this.password }));
            window.location.href = '/'
          })
          .catch(error => {
            alert(error.message);
          });
        },
      },

    template: `
    <div class="container">
        <form @submit.prevent="login">
            <input v-model="password" type="password" placeholder="Введите пароль" />
            <button type="submit">Войти</button>
        </form>
    </div>
    `
}