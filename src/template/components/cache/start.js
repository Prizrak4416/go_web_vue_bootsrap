export default {
    data() {
        return {
            message: "Загрузка данных..."
        }
    },
    methods: {
        fetchData() {
            fetch('http://localhost:8080/api/data')
                .then(response => response.json())
                .then(data => {
                    this.message = data.message;
                })
                .catch(error => {
                    this.message = 'Ошибка при получении данных';
                    console.error('Ошибка при выполнении fetch запроса:', error);
                });
        }
    },
    mounted() {
        this.fetchData();
    },
    template: `<div>{{ message }}</div>`
}