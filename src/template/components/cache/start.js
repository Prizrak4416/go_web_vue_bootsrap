export default {
    data() {
        return {
            message: "Загрузка данных...",
            data: [],
            message_post : '',
            response: ''
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
        }
    },
    mounted() {
        this.fetchData()
    },
    template: `
    <div>
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
        
    </div>
    `
}