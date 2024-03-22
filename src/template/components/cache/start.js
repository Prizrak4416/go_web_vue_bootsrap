export default {
    data() {
        return {
            message: "Загрузка данных...",
            message_ws: '',
            data: [],
            message_post : '',
            response: '',
            ws: null, // WebSocket connection
        }
    },
    methods: {
        fetchData() {
            fetch('/api/data')
                .then(response => response.json())
                .then(data => {
                    this.data = data.message
                    this.message = ''
                })
                .catch(error => {
                    this.message = 'Ошибка при получении данных'
                    console.error('Ошибка при выполнении fetch запроса:', error)
                });
        },
        sendUserData() {
            const userData = {
                role: "user",
                message_post: "Данные пользователя:"
            };
            this.sendData(userData);
        },
        GetUptime() {
            const userData = {
                role: "getut",
                message_post: "Получение времени:"
            };
            this.sendData(userData);
        },
        sendAdminData() {
            const adminData = {
                role: "admin",
                message_post: "Данные администратора:"
            };
            this.sendData(adminData);
        },
        getSSH() {
            const adminData = {
                role: "getssh",
                message_post: "Получение SSH ключей:"
            };
            this.sendData(adminData);
        },
        sendData(data) {
            fetch('/api/data', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(data => {
                this.response = data.response;
            })
            .catch(error => {
                this.message = 'Ошибка при отправке данных'
                console.error('Ошибка при выполнении fetch запроса:', error)
            });
        },
        sendWebs() {
            this.ws.send('Првет сервер')
        },
        sendWebs2() {
            this.ws.send('Еще 1 привет')
        }

    },
    mounted() {
        this.fetchData()
        // this.ws = new WebSocket("ws");
        this.ws = new WebSocket(((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/ws");

        this.ws.onmessage = (event) => {
            // При получении сообщения обновляем данные
            this.message_ws = event.data;
        };

        this.ws.onopen = () => {
            // WebSocket соединение открыто
            console.log("WebSocket connection is open.");
        };

        this.ws.onerror = (error) => {
            // Обработка ошибок WebSocket соединения
            console.error("WebSocket error:", error);
        };

        this.ws.onclose = () => {
            // WebSocket соединение закрыто
            console.log("WebSocket connection is closed.");
        };
    },

    beforeUnmount() {
        // Закрываем WebSocket соединение, когда компонент удаляется
        if (this.ws) {
            this.ws.close();
        }
    },


    template: `
    <div>
        <p>{{ message_ws }}</p>
        <hr>
        <br>
        {{ message }}
        <div v-for="value in data">
            <strong>ID: {{ value['ID'] }}</strong> |&nbsp
            <b>Name: {{ value['Name'] }}</b> |&nbsp
            <b>UserName: {{ value['UserName'] }}</b>
            <br>
        </div>
        <hr>

        <div>
            <button @click="sendUserData">Отправить данные пользователя</button>&nbsp
            <button @click="sendAdminData">Отправить данные администратора</button>&nbsp
            <button @click="GetUptime">Получить время работы Linux</button>&nbsp
            <button @click="getSSH">Получить ssh</button>
            <div v-if="typeof response === 'string'">{{ response }}</div>
            <div v-else>
                <p v-for="value in response">
                    {{ value }}
                </p>
            </div>
            <br>
            {{ typeof response }}
        </div>
        <hr>
        <br>
        <div>
            <button @click="sendWebs">send socket</button>&nbsp
            <button @click="sendWebs2">send socket 2</button>&nbsp
        </div>
    </div>
    `
}