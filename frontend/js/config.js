window.APP_CONFIG = {
    BASE_URL: "http://core-service:8099",  // внутри docker сети
    ENDPOINTS: {
        login: "/users/login",
        register: "/users",
        listMedicalData: (userId) => `/medical-data?user_id=${userId}`,
        deleteMedicalData: (id) => `/medical-data/${id}`,
    },
};
