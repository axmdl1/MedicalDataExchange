// Single Page Application
class App {
    constructor() {
        this.currentPage = null;
        this.appContainer = document.getElementById('app');
        this.userNameCache = new Map(); // Кэш: user_id → "Имя Фамилия"
        this.init();
        this.initTheme();
    }

    init() {
        // Check authentication
        if (!auth.isAuthenticated()) {
            this.renderLogin();
        } else {
            this.renderDashboard();
        }

        // Handle navigation
        window.addEventListener('hashchange', () => this.handleRoute());
        this.handleRoute();
    }

    initTheme() {
        // Load saved theme or default to light
        const savedTheme = localStorage.getItem('theme') || 'light';
        document.documentElement.setAttribute('data-theme', savedTheme);
    }

    toggleTheme() {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
    }

    handleRoute() {
        const hash = window.location.hash.slice(1) || 'dashboard';

        if (!auth.isAuthenticated() && hash !== 'login') {
            window.location.hash = 'login';
            return;
        }

        switch (hash) {
            case 'login':
                this.renderLogin();
                break;
            case 'dashboard':
                this.renderDashboard();
                break;
            case 'patients':
                this.renderPatients();
                break;
            case 'employees':
                this.renderEmployees();
                break;
            case 'medical-data':
                this.renderMedicalData();
                break;
            case 'transfers':
                this.renderTransfers();
                break;
            case 'clinics':
                this.renderClinics();
                break;
            case 'event-log':
                this.renderEventLog();
                break;
            case 'create-user':
                this.renderCreateUser();
                break;
            case 'create-medical-data':
                this.renderCreateMedicalData();
                break;
            case 'create-transfer':
                this.renderCreateTransfer();
                break;
            default:
                this.renderDashboard();
        }
    }

    // Render sidebar navigation
    renderSidebar() {
        const role = auth.getRole();

        const navItems = [];

        // Common for all
        navItems.push({ icon: 'fa-chart-line', label: 'Dashboard', hash: 'dashboard' });

        if (role === ROLES.PATIENT) {
            navItems.push({ icon: 'fa-file-medical', label: 'My Medical Data', hash: 'medical-data' });
            navItems.push({ icon: 'fa-exchange-alt', label: 'Transfer Requests', hash: 'transfers' });
            navItems.push({ icon: 'fa-history', label: 'Event Log', hash: 'event-log' });
        } else if (role === ROLES.EMPLOYEE) {
            navItems.push({ icon: 'fa-users', label: 'Patients', hash: 'patients' });
            navItems.push({ icon: 'fa-user-tie', label: 'Employees', hash: 'employees' });
            navItems.push({ icon: 'fa-file-medical', label: 'Medical Records', hash: 'medical-data' });
            navItems.push({ icon: 'fa-exchange-alt', label: 'Data Transfers', hash: 'transfers' });
            navItems.push({ icon: 'fa-hospital', label: 'Clinics', hash: 'clinics' });
            navItems.push({ icon: 'fa-history', label: 'Event Log', hash: 'event-log' });
            navItems.push({ icon: 'fa-user-plus', label: 'Create User', hash: 'create-user' });
            navItems.push({ icon: 'fa-file-medical-alt', label: 'Create Medical Data', hash: 'create-medical-data' });
            navItems.push({ icon: 'fa-paper-plane', label: 'Create Transfer', hash: 'create-transfer' });
        } else if (role === ROLES.ADMIN) {
            navItems.push({ icon: 'fa-users', label: 'All Patients', hash: 'patients' });
            navItems.push({ icon: 'fa-user-tie', label: 'All Employees', hash: 'employees' });
            navItems.push({ icon: 'fa-file-medical', label: 'All Medical Data', hash: 'medical-data' });
            navItems.push({ icon: 'fa-exchange-alt', label: 'All Transfers', hash: 'transfers' });
            navItems.push({ icon: 'fa-hospital', label: 'Clinics', hash: 'clinics' });
            navItems.push({ icon: 'fa-history', label: 'Event Log', hash: 'event-log' });
            navItems.push({ icon: 'fa-user-plus', label: 'Create User', hash: 'create-user' });
            navItems.push({ icon: 'fa-file-medical-alt', label: 'Create Medical Data', hash: 'create-medical-data' });
            navItems.push({ icon: 'fa-paper-plane', label: 'Create Transfer', hash: 'create-transfer' });
        }

        const navHTML = navItems.map(item => `
            <a href="#${item.hash}" class="nav-item ${window.location.hash.slice(1) === item.hash ? 'active' : ''}">
                <i class="fas ${item.icon}"></i> ${item.label}
            </a>
        `).join('');

        return `
            <aside class="sidebar">
                <div class="logo">
                    <div class="logo-icon">
                        <i class="fas fa-shield-alt"></i>
                        <span style="font-size: 18px;">Data Transfer System</span>
                    </div>
                </div>
                <nav class="sidebar-nav">
                    ${navHTML}
                    <a href="#" onclick="app.logout(); return false;" class="nav-item">
                        <i class="fas fa-sign-out-alt"></i> Logout
                    </a>
                </nav>
                <div class="theme-toggle">
                    <button onclick="app.toggleTheme()" title="Toggle Theme">
                        <i class="fas fa-moon"></i>
                    </button>
                </div>
            </aside>
        `;
    }

    renderLogin() {
        this.appContainer.innerHTML = `
            <div class="login-container">
                <div class="login-box">
                    <h1>Medical Data System</h1>

                    <div class="forms-wrapper">
                        <!-- Login Form -->
                        <div class="form-section">
                            <h2>Login</h2>
                            <form id="loginForm">
                                <div class="form-group">
                                    <input type="email" id="loginEmail" placeholder="Email address" required>
                                </div>
                                <div class="form-group">
                                    <input type="password" id="loginPassword" placeholder="Password" required>
                                </div>
                                <div id="loginError" style="color: red; font-size: 14px; margin-bottom: 10px;"></div>
                                <button type="submit" class="btn btn-primary">Sign in</button>
                            </form>
                        </div>

                        <!-- Register Form -->
                        <div class="form-section">
                            <h2>Register as Patient</h2>
                            <form id="registerForm">
                                <div class="form-group">
                                    <input type="text" id="regFirstName" placeholder="First Name" required>
                                </div>
                                <div class="form-group">
                                    <input type="text" id="regLastName" placeholder="Last Name" required>
                                </div>
                                <div class="form-group">
                                    <input type="email" id="regEmail" placeholder="Email address" required>
                                </div>
                                <div class="form-group">
                                    <input type="tel" id="regPhone" placeholder="Phone Number" required>
                                </div>
                                <div class="form-group">
                                    <input type="password" id="regPassword" placeholder="Password" required>
                                </div>
                                <div id="registerError" style="color: red; font-size: 14px; margin-bottom: 10px;"></div>
                                <button type="submit" class="btn btn-secondary">Sign up</button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        `;

        // Login form handler
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('loginEmail').value;
            const password = document.getElementById('loginPassword').value;

            try {
                const response = await api.login(email, password);
                // Backend возвращает access_token вместо token
                auth.saveAuth(response.access_token, response.user);
                window.location.hash = 'dashboard';
            } catch (error) {
                document.getElementById('loginError').textContent = error.message || 'Login failed';
            }
        });

        // Register form handler
        document.getElementById('registerForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const userData = {
                first_name: document.getElementById('regFirstName').value,
                last_name: document.getElementById('regLastName').value,
                email: document.getElementById('regEmail').value,
                phone_number: document.getElementById('regPhone').value,
                password: document.getElementById('regPassword').value,
                type: ROLES.PATIENT
            };

            try {
                await api.createUser(userData);
                // Auto login after registration
                const response = await api.login(userData.email, userData.password);
                // Backend возвращает access_token вместо token
                auth.saveAuth(response.access_token, response.user);
                window.location.hash = 'dashboard';
            } catch (error) {
                document.getElementById('registerError').textContent = error.message || 'Registration failed';
            }
        });
    }

    async renderDashboard() {
        const role = auth.getRole();
        const user = auth.getUser();

        let stats = {};
        let chartData = {};

        try {
            if (role === ROLES.ADMIN) {
                // Admin sees all data
                const [users, medicalData, transfers] = await Promise.all([
                    api.getUsers(),
                    api.getMedicalData(),
                    api.getTransfers()
                ]);

                const allUsers = users.users || [];
                const allMedicalData = medicalData.medical_data || [];
                const allTransfers = transfers.transfers || [];

                stats = {
                    patients: allUsers.filter(u => u.type === ROLES.PATIENT).length,
                    employees: allUsers.filter(u => u.type === ROLES.EMPLOYEE).length,
                    records: allMedicalData.length,
                    transfers: allTransfers.length
                };

                // Chart data for admin
                chartData = {
                    usersByType: {
                        patients: stats.patients,
                        employees: stats.employees,
                        admins: allUsers.filter(u => u.type === ROLES.ADMIN).length
                    },
                    transfersByStatus: {
                        pending: allTransfers.filter(t => t.status === TRANSFER_STATUS.PENDING).length,
                        confirmed: allTransfers.filter(t => t.status === TRANSFER_STATUS.CONFIRMED).length,
                        rejected: allTransfers.filter(t => t.status === TRANSFER_STATUS.REJECTED).length
                    }
                };
            } else if (role === ROLES.EMPLOYEE) {
                // Employee sees clinic data
                const clinicId = auth.getClinicId();
                const [medicalData, transfers] = await Promise.all([
                    api.getMedicalData({ clinic_id: clinicId }),
                    api.getTransfers({ clinic_id: clinicId })
                ]);

                const allMedicalData = medicalData.medical_data || [];
                const allTransfers = transfers.transfers || [];

                stats = {
                    patients: allMedicalData.length,
                    records: allMedicalData.length,
                    transfers: allTransfers.length
                };

                // Chart data for employee
                chartData = {
                    transfersByStatus: {
                        pending: allTransfers.filter(t => t.status === TRANSFER_STATUS.PENDING).length,
                        confirmed: allTransfers.filter(t => t.status === TRANSFER_STATUS.CONFIRMED).length,
                        rejected: allTransfers.filter(t => t.status === TRANSFER_STATUS.REJECTED).length
                    }
                };
            } else {
                // Patient sees own data
                const userId = user.id;
                const [medicalData, transfers] = await Promise.all([
                    api.getMedicalData({ user_id: userId }),
                    api.getTransfers({ user_id: userId })
                ]);

                const allTransfers = transfers.transfers || [];

                stats = {
                    records: medicalData.medical_data?.length || 0,
                    pending: allTransfers.filter(t => t.status === TRANSFER_STATUS.PENDING).length,
                    approved: allTransfers.filter(t => t.status === TRANSFER_STATUS.CONFIRMED).length
                };

                // Chart data for patient
                chartData = {
                    transfersByStatus: {
                        pending: stats.pending,
                        confirmed: stats.approved,
                        rejected: allTransfers.filter(t => t.status === TRANSFER_STATUS.REJECTED).length
                    }
                };
            }
        } catch (error) {
            console.error('Failed to load dashboard data:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <h1 style="margin-bottom: 20px;">Welcome, ${user.first_name} ${user.last_name}</h1>
                    <p style="color: #666; margin-bottom: 30px;">Role: ${role.toUpperCase()}</p>

                    <div class="stats-grid">
                        ${role === ROLES.PATIENT ? `
                            <div class="stat-card">
                                <div class="stat-header"><span>My Medical Records</span></div>
                                <div class="stat-value">${stats.records}</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-header"><span>Pending Requests</span></div>
                                <div class="stat-value">${stats.pending}</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-header"><span>Approved Transfers</span></div>
                                <div class="stat-value">${stats.approved}</div>
                            </div>
                        ` : `
                            <div class="stat-card">
                                <div class="stat-header"><span>Patients</span></div>
                                <div class="stat-value">${stats.patients}</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-header"><span>Medical Records</span></div>
                                <div class="stat-value">${stats.records}</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-header"><span>Transfers</span></div>
                                <div class="stat-value">${stats.transfers}</div>
                            </div>
                        `}
                    </div>

                    <div class="widget" style="margin-top: 30px;">
                        <div class="widget-header">
                            <h3>Quick Actions</h3>
                        </div>
                        <div style="display: flex; gap: 10px; margin-top: 15px;">
                            ${role === ROLES.EMPLOYEE || role === ROLES.ADMIN ? `
                                <button class="btn btn-primary btn-sm" onclick="window.location.hash='patients'">View Patients</button>
                                <button class="btn btn-secondary btn-sm" onclick="window.location.hash='medical-data'">View Medical Data</button>
                            ` : `
                                <button class="btn btn-primary btn-sm" onclick="window.location.hash='medical-data'">View My Records</button>
                                <button class="btn btn-secondary btn-sm" onclick="window.location.hash='transfers'">Transfer Requests</button>
                            `}
                        </div>
                    </div>

                    ${role === ROLES.ADMIN || role === ROLES.EMPLOYEE ? `
                        <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 20px; margin-top: 30px;">
                            ${role === ROLES.ADMIN && chartData.usersByType ? `
                                <div class="widget">
                                    <div class="widget-header">
                                        <h3>Users by Type</h3>
                                    </div>
                                    <div style="padding: 20px;">
                                        <canvas id="usersChart"></canvas>
                                    </div>
                                </div>
                            ` : ''}

                            ${chartData.transfersByStatus ? `
                                <div class="widget">
                                    <div class="widget-header">
                                        <h3>Transfer Requests Status</h3>
                                    </div>
                                    <div style="padding: 20px;">
                                        <canvas id="transfersChart"></canvas>
                                    </div>
                                </div>
                            ` : ''}
                        </div>
                    ` : ''}

                    ${role === ROLES.PATIENT && chartData.transfersByStatus ? `
                        <div class="widget" style="margin-top: 30px;">
                            <div class="widget-header">
                                <h3>My Transfer Requests</h3>
                            </div>
                            <div style="padding: 20px; max-width: 500px; margin: 0 auto;">
                                <canvas id="transfersChart"></canvas>
                            </div>
                        </div>
                    ` : ''}
                </main>
            </div>
        `;

        // Render charts after DOM is ready
        setTimeout(() => {
            if (role === ROLES.ADMIN && chartData.usersByType) {
                this.renderUsersChart(chartData.usersByType);
            }
            if (chartData.transfersByStatus) {
                this.renderTransfersChart(chartData.transfersByStatus);
            }
        }, 100);
    }

    renderUsersChart(data) {
        const ctx = document.getElementById('usersChart');
        if (!ctx) return;

        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: ['Patients', 'Employees', 'Admins'],
                datasets: [{
                    data: [data.patients, data.employees, data.admins],
                    backgroundColor: [
                        '#4CAF50',
                        '#2196F3',
                        '#FF9800'
                    ],
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: {
                            padding: 15,
                            font: {
                                size: 12
                            }
                        }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = ((value / total) * 100).toFixed(1);
                                return `${label}: ${value} (${percentage}%)`;
                            }
                        }
                    }
                }
            }
        });
    }

    renderTransfersChart(data) {
        const ctx = document.getElementById('transfersChart');
        if (!ctx) return;

        new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['Pending', 'Confirmed', 'Rejected'],
                datasets: [{
                    label: 'Transfer Requests',
                    data: [data.pending, data.confirmed, data.rejected],
                    backgroundColor: [
                        '#FFC107',
                        '#4CAF50',
                        '#F44336'
                    ],
                    borderWidth: 0,
                    borderRadius: 6
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: true,
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            stepSize: 1,
                            font: {
                                size: 11
                            }
                        },
                        grid: {
                            color: '#f0f0f0'
                        }
                    },
                    x: {
                        grid: {
                            display: false
                        },
                        ticks: {
                            font: {
                                size: 11
                            }
                        }
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `Count: ${context.parsed.y}`;
                            }
                        }
                    }
                }
            }
        });
    }

    async renderPatients() {
        const role = auth.getRole();
        let users = [];

        try {
            const response = await api.getUsers();
            users = response.users || [];

            // Filter only patients
            users = users.filter(u => u.type === ROLES.PATIENT);

            // Employee sees only his clinic patients (but patient has no clinic_id!)
            // So we need to check which clinic has medical records for this patient
            if (role === ROLES.EMPLOYEE) {
                const myClinicId = auth.getClinicId();
                // For now, show all patients - filtering should be done on backend
                // Or we fetch medical_data and see which patients belong to this clinic
            }
        } catch (error) {
            console.error('Failed to load patients:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px;">
                        <h1 style="margin: 0;">${role === ROLES.ADMIN ? 'All Patients' : 'Patients'}</h1>
                        <button class="btn btn-primary" onclick="app.showUserModal('patient')"><i class="fas fa-plus"></i> Create Patient</button>
                    </div>

                    <div class="widget">
                        <table class="data-table">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Email</th>
                                    <th>Phone</th>
                                    <th>Type</th>
                                    <th>Created</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${users.length === 0 ? `
                                    <tr>
                                        <td colspan="6" style="text-align: center; padding: 20px; color: var(--text-secondary);">
                                            No patients found
                                        </td>
                                    </tr>
                                ` : users.map(user => `
                                    <tr>
                                        <td>${user.first_name} ${user.last_name}</td>
                                        <td>${user.email}</td>
                                        <td>${user.phone_number || 'N/A'}</td>
                                        <td><span class="badge badge-success">${user.type}</span></td>
                                        <td>${new Date(user.created_at).toLocaleDateString()}</td>
                                        <td>
                                            <button class="btn btn-primary btn-sm" onclick="app.showUserModal('patient', ${user.id})"><i class="fas fa-edit"></i> Edit</button>
                                            ${role === ROLES.ADMIN ? `<button class="btn btn-danger btn-sm" onclick="app.deleteUser(${user.id})"><i class="fas fa-trash"></i> Delete</button>` : ''}
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </main>
            </div>

            <!-- User Modal -->
            <div id="userModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); z-index: 2000; align-items: center; justify-content: center;">
                <div style="background: var(--bg-primary); padding: 32px; border-radius: 12px; max-width: 600px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 50px var(--shadow);">
                    <h2 id="userModalTitle" style="margin-bottom: 24px; color: var(--text-primary);">Create Patient</h2>
                    <form id="userForm">
                        <input type="hidden" id="userId">
                        <input type="hidden" id="userType">
                        <div class="form-group">
                            <label>First Name</label>
                            <input type="text" id="userFirstName" required>
                        </div>
                        <div class="form-group">
                            <label>Last Name</label>
                            <input type="text" id="userLastName" required>
                        </div>
                        <div class="form-group">
                            <label>Email</label>
                            <input type="email" id="userEmail" required>
                        </div>
                        <div class="form-group">
                            <label>Phone</label>
                            <input type="tel" id="userPhone" required>
                        </div>
                        <div class="form-group" id="userPasswordGroup">
                            <label>Password</label>
                            <input type="password" id="userPassword" placeholder="Leave empty to keep current">
                        </div>
                        <div class="form-group" id="userClinicGroup" style="display: none;">
                            <label>Clinic ID</label>
                            <input type="number" id="userClinicId">
                        </div>
                        <div class="form-actions">
                            <button type="submit" class="btn btn-primary"><i class="fas fa-save"></i> Save</button>
                            <button type="button" class="btn btn-secondary" onclick="app.hideUserModal()"><i class="fas fa-times"></i> Cancel</button>
                        </div>
                    </form>
                </div>
            </div>
        `;

        // Setup form handler
        document.getElementById('userForm').addEventListener('submit', (e) => this.submitUserForm(e));
    }

    async renderMedicalData() {
        const role = auth.getRole();
        const user = auth.getUser();
        let medicalData = [];

        try {
            const filters = {};
            if (role === ROLES.PATIENT) {
                filters.user_id = user.id;
            } else if (role === ROLES.EMPLOYEE) {
                filters.clinic_id = auth.getClinicId();
            }

            const response = await api.getMedicalData(filters);
            medicalData = response.medical_data || [];

            // Предзагружаем имена всех пациентов
            const userIds = [...new Set(medicalData.map(r => r.user_id))];
            await Promise.all(userIds.map(id => this.getUserName(id)));

        } catch (error) {
            console.error('Failed to load medical data:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px;">
                        <h1 style="margin: 0;">Medical Records</h1>
                        ${role !== ROLES.PATIENT ? '<button class="btn btn-primary" onclick="app.showMedicalDataModal()"><i class="fas fa-plus"></i> Create Record</button>' : ''}
                    </div>

                    <div class="widget">
                        <table class="data-table">
                            <thead>
                                <tr>
                                    <th>Patient</th>
                                    <th>Date</th>
                                    <th>Diagnosis</th>
                                    <th>Complaint</th>
                                    <th>Treatment</th>
                                    <th>Medications</th>
                                    ${role !== ROLES.PATIENT ? '<th>Actions</th>' : ''}
                                </tr>
                            </thead>
                            <tbody>
                                ${medicalData.length === 0 ? `
                                    <tr>
                                        <td colspan="${role !== ROLES.PATIENT ? '7' : '6'}" style="text-align: center; padding: 20px; color: var(--text-secondary);">
                                            No records found
                                        </td>
                                    </tr>
                                ` : medicalData.map(record => `
                                    <tr>
                                        <td>${this.escapeHtml(this.userNameCache.get(record.user_id) || 'Loading...')}</td>
                                        <td>${new Date(record.created_at).toLocaleDateString()}</td>
                                        <td>${this.escapeHtml(record.diagnosis || 'N/A')}</td>
                                        <td>${this.escapeHtml(record.complaint || 'N/A')}</td>
                                        <td>${this.escapeHtml(record.treatment || 'N/A')}</td>
                                        <td>${this.escapeHtml(record.medications || 'N/A')}</td>
                                        ${role !== ROLES.PATIENT ? `
                                            <td>
                                                <button class="btn btn-danger btn-sm" onclick="app.deleteMedicalData(${record.id})"><i class="fas fa-trash"></i> Delete</button>
                                            </td>
                                        ` : ''}
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </main>
            </div>

            <!-- Medical Data Modal -->
            <div id="medicalDataModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); z-index: 2000; align-items: center; justify-content: center;">
                <div style="background: var(--bg-primary); padding: 32px; border-radius: 12px; max-width: 600px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 50px var(--shadow);">
                    <h2 id="medicalDataModalTitle" style="margin-bottom: 24px; color: var(--text-primary);">Create Medical Record</h2>
                    <form id="medicalDataForm">
                        <div class="form-group">
                            <label>Patient</label>
                            <select id="medicalDataPatientSelect" required>
                                <option value="">Select patient...</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label>Clinic</label>
                            <select id="medicalDataClinicSelect" required>
                                <option value="">Select clinic...</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label>Diagnosis</label>
                            <input type="text" id="medicalDataDiagnosis" required>
                        </div>
                        <div class="form-group">
                            <label>Complaint</label>
                            <textarea id="medicalDataComplaint" rows="3" required></textarea>
                        </div>
                        <div class="form-group">
                            <label>Treatment</label>
                            <textarea id="medicalDataTreatment" rows="3" required></textarea>
                        </div>
                        <div class="form-group">
                            <label>Medications</label>
                            <textarea id="medicalDataMedications" rows="3" required></textarea>
                        </div>
                        <div class="form-actions">
                            <button type="submit" class="btn btn-primary"><i class="fas fa-save"></i> Save</button>
                            <button type="button" class="btn btn-secondary" onclick="app.hideMedicalDataModal()"><i class="fas fa-times"></i> Cancel</button>
                        </div>
                    </form>
                </div>
            </div>
        `;

        // Setup form handler
        if (role !== ROLES.PATIENT) {
            document.getElementById('medicalDataForm').addEventListener('submit', (e) => this.submitMedicalDataForm(e));
        }
    }

    async renderTransfers() {
        const role = auth.getRole();
        const user = auth.getUser();
        let transfers = [];

        try {
            const filters = {};
            if (role === ROLES.PATIENT) {
                filters.user_id = user.id;
            } else if (role === ROLES.EMPLOYEE) {
                filters.clinic_id = auth.getClinicId();
            }

            const response = await api.getTransfers(filters);
            transfers = response.transfers || [];

            // Предзагружаем имена
            const userIds = [...new Set(transfers.map(t => t.user_id))];
            await Promise.all(userIds.map(id => this.getUserName(id)));

        } catch (error) {
            console.error('Failed to load transfers:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <h1 style="margin-bottom: 30px;">Data Transfer Requests</h1>

                    <div class="widget">
                        <table class="data-table">
                            <thead>
                                <tr>
                                    <th>Patient</th>
                                    <th>Date</th>
                                    <th>From Clinic</th>
                                    <th>To Clinic</th>
                                    <th>Status</th>
                                    ${role === ROLES.PATIENT ? '<th>Action</th>' : ''}
                                </tr>
                            </thead>
                            <tbody>
                                ${transfers.length === 0 ? `
                                    <tr>
                                        <td colspan="${role === ROLES.PATIENT ? '6' : '5'}" style="text-align: center; padding: 20px; color: #666;">
                                            No transfers found
                                        </td>
                                    </tr>
                                ` : transfers.map(transfer => `
                                    <tr>
                                        <td>${this.escapeHtml(this.userNameCache.get(transfer.user_id) || 'Unknown')}</td>
                                        <td>${new Date(transfer.created_at).toLocaleDateString()}</td>
                                        <td>Clinic #${transfer.from_clinic_id}</td>
                                        <td>Clinic #${transfer.to_clinic_id}</td>
                                        <td>
                                            <span class="badge badge-${transfer.status === 'confirmed' ? 'success' : transfer.status === 'pending' ? 'warning' : 'inactive'}">
                                                ${transfer.status}
                                            </span>
                                        </td>
                                        ${role === ROLES.PATIENT && transfer.status === 'pending' ? `
                                            <td>
                                                <button class="btn btn-primary btn-sm" onclick="app.approveTransfer(${transfer.id})">Approve</button>
                                                <button class="btn btn-danger btn-sm" onclick="app.rejectTransfer(${transfer.id})">Reject</button>
                                            </td>
                                        ` : role === ROLES.PATIENT ? '<td>-</td>' : ''}
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </main>
            </div>
        `;
    }

    async approveTransfer(transferId) {
        try {
            await api.makeTransferDecision(transferId, true);
            alert('Transfer approved successfully!');
            this.renderTransfers();
        } catch (error) {
            alert('Failed to approve transfer: ' + error.message);
        }
    }

    async rejectTransfer(transferId) {
        try {
            await api.makeTransferDecision(transferId, false);
            alert('Transfer rejected successfully!');
            this.renderTransfers();
        } catch (error) {
            alert('Failed to reject transfer: ' + error.message);
        }
    }

    async renderEmployees() {
        const role = auth.getRole();
        let users = [];

        try {
            const response = await api.getUsers();
            console.log('getUsers response:', response);
            users = response.users || [];
            console.log('All users:', users);

            // Filter only employees
            users = users.filter(u => u.type === ROLES.EMPLOYEE);
            console.log('Filtered employees:', users);

            // Employee sees only his clinic
            if (role === ROLES.EMPLOYEE) {
                const myClinicId = auth.getClinicId();
                console.log('My clinic ID:', myClinicId);
                users = users.filter(u => {
                    const userClinicId = u.clinic_id?.value || u.clinic_id;
                    console.log('User clinic ID:', userClinicId, 'Match:', userClinicId == myClinicId);
                    return userClinicId == myClinicId; // Use loose equality to handle type conversion
                });
                console.log('Filtered by clinic:', users);
            }
        } catch (error) {
            console.error('Failed to load employees:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px;">
                        <h1 style="margin: 0;">${role === ROLES.ADMIN ? 'All Employees' : 'Clinic Employees'}</h1>
                        ${role === ROLES.ADMIN ? '<button class="btn btn-primary" onclick="app.showUserModal(\'employee\')"><i class="fas fa-plus"></i> Create Employee</button>' : ''}
                    </div>

                    <div class="widget">
                        <table class="data-table">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Email</th>
                                    <th>Phone</th>
                                    <th>Clinic ID</th>
                                    <th>Created</th>
                                    ${role === ROLES.ADMIN ? '<th>Actions</th>' : ''}
                                </tr>
                            </thead>
                            <tbody>
                                ${users.length === 0 ? `
                                    <tr>
                                        <td colspan="${role === ROLES.ADMIN ? '6' : '5'}" style="text-align: center; padding: 20px; color: var(--text-secondary);">
                                            No employees found
                                        </td>
                                    </tr>
                                ` : users.map(user => `
                                    <tr>
                                        <td>${user.first_name} ${user.last_name}</td>
                                        <td>${user.email}</td>
                                        <td>${user.phone_number || 'N/A'}</td>
                                        <td>Clinic #${user.clinic_id?.value || user.clinic_id}</td>
                                        <td>${new Date(user.created_at).toLocaleDateString()}</td>
                                        ${role === ROLES.ADMIN ? `
                                            <td>
                                                <button class="btn btn-primary btn-sm" onclick="app.showUserModal('employee', ${user.id})"><i class="fas fa-edit"></i> Edit</button>
                                                <button class="btn btn-danger btn-sm" onclick="app.deleteUser(${user.id})"><i class="fas fa-trash"></i> Delete</button>
                                            </td>
                                        ` : ''}
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </main>
            </div>

            <!-- User Modal -->
            <div id="userModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); z-index: 2000; align-items: center; justify-content: center;">
                <div style="background: var(--bg-primary); padding: 32px; border-radius: 12px; max-width: 600px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 50px var(--shadow);">
                    <h2 id="userModalTitle" style="margin-bottom: 24px; color: var(--text-primary);">Create Employee</h2>
                    <form id="userForm">
                        <input type="hidden" id="userId">
                        <input type="hidden" id="userType">
                        <div class="form-group">
                            <label>First Name</label>
                            <input type="text" id="userFirstName" required>
                        </div>
                        <div class="form-group">
                            <label>Last Name</label>
                            <input type="text" id="userLastName" required>
                        </div>
                        <div class="form-group">
                            <label>Email</label>
                            <input type="email" id="userEmail" required>
                        </div>
                        <div class="form-group">
                            <label>Phone</label>
                            <input type="tel" id="userPhone" required>
                        </div>
                        <div class="form-group" id="userPasswordGroup">
                            <label>Password</label>
                            <input type="password" id="userPassword" placeholder="Leave empty to keep current">
                        </div>
                        <div class="form-group" id="userClinicGroup">
                            <label>Clinic ID</label>
                            <input type="number" id="userClinicId" required>
                        </div>
                        <div class="form-actions">
                            <button type="submit" class="btn btn-primary"><i class="fas fa-save"></i> Save</button>
                            <button type="button" class="btn btn-secondary" onclick="app.hideUserModal()"><i class="fas fa-times"></i> Cancel</button>
                        </div>
                    </form>
                </div>
            </div>
        `;

        // Setup form handler
        if (role === ROLES.ADMIN) {
            document.getElementById('userForm').addEventListener('submit', (e) => this.submitUserForm(e));
        }
    }

    async renderCreateUser() {
        const role = auth.getRole();
        const myClinicId = auth.getClinicId();

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <h1 style="margin-bottom: 30px;">Create New User</h1>

                    <div class="widget" style="max-width: 600px;">
                        <form id="createUserForm">
                            <div class="form-group">
                                <label>User Type</label>
                                <select id="userType" required>
                                    <option value="patient">Patient</option>
                                    ${role === ROLES.ADMIN ? '<option value="employee">Employee</option>' : ''}
                                    ${role === ROLES.ADMIN ? '<option value="admin">Admin</option>' : ''}
                                </select>
                            </div>

                            <div class="form-group">
                                <label>First Name</label>
                                <input type="text" id="firstName" required>
                            </div>

                            <div class="form-group">
                                <label>Last Name</label>
                                <input type="text" id="lastName" required>
                            </div>

                            <div class="form-group">
                                <label>Email</label>
                                <input type="email" id="email" required>
                            </div>

                            <div class="form-group">
                                <label>Phone</label>
                                <input type="tel" id="phone" required>
                            </div>

                            <div class="form-group">
                                <label>Password</label>
                                <input type="password" id="password" value="password123" required>
                            </div>

                            <div class="form-group" id="clinicIdGroup" style="display: none;">
                                <label>Clinic ID</label>
                                <input type="number" id="clinicId" value="${myClinicId || 1}">
                            </div>

                            <div id="createUserError" style="color: red; margin: 10px 0;"></div>
                            <div id="createUserSuccess" style="color: green; margin: 10px 0;"></div>

                            <div class="form-actions">
                                <button type="submit" class="btn btn-primary">Create User</button>
                                <button type="button" class="btn btn-secondary" onclick="window.location.hash='dashboard'">Cancel</button>
                            </div>
                        </form>
                    </div>
                </main>
            </div>
        `;

        // Show clinic field for employee type
        document.getElementById('userType').addEventListener('change', (e) => {
            const clinicGroup = document.getElementById('clinicIdGroup');
            if (e.target.value === 'employee') {
                clinicGroup.style.display = 'block';
                if (role === ROLES.EMPLOYEE) {
                    document.getElementById('clinicId').value = myClinicId;
                    document.getElementById('clinicId').readOnly = true;
                }
            } else {
                clinicGroup.style.display = 'none';
            }
        });

        // Form submit
        document.getElementById('createUserForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const userType = document.getElementById('userType').value;
            const userData = {
                first_name: document.getElementById('firstName').value,
                last_name: document.getElementById('lastName').value,
                email: document.getElementById('email').value,
                phone_number: document.getElementById('phone').value,
                password: document.getElementById('password').value,
                type: userType
            };

            // For employee type, set clinic_id
            if (userType === 'employee') {
                userData.clinic_id = parseInt(document.getElementById('clinicId').value);
            }
            // For patient type created by employee, also set clinic_id
            else if (userType === 'patient' && role === ROLES.EMPLOYEE) {
                userData.clinic_id = myClinicId;
            }

            try {
                await api.createUser(userData);
                document.getElementById('createUserSuccess').textContent = 'User created successfully!';
                document.getElementById('createUserError').textContent = '';
                document.getElementById('createUserForm').reset();
            } catch (error) {
                document.getElementById('createUserError').textContent = error.message || 'Failed to create user';
                document.getElementById('createUserSuccess').textContent = '';
            }
        });
    }

    async renderCreateMedicalData() {
        const role = auth.getRole();
        const myClinicId = auth.getClinicId();

        // Load patients and clinics for selection
        let patients = [];
        let clinics = [];
        try {
            const response = await api.getUsers();
            patients = (response.users || []).filter(u => u.type === ROLES.PATIENT);

            const clinicsResponse = await api.getClinics();
            clinics = clinicsResponse.clinics || [];
        } catch (error) {
            console.error('Failed to load data:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <h1 style="margin-bottom: 30px;">Create Medical Record</h1>

                    <div class="widget" style="max-width: 800px;">
                        <form id="createMedicalForm">
                            <div class="form-row">
                                <div class="form-group">
                                    <label>Patient</label>
                                    <select id="patientId" required>
                                        <option value="">Select patient...</option>
                                        ${patients.map(p => `
                                            <option value="${p.id}">${p.first_name} ${p.last_name} (${p.email})</option>
                                        `).join('')}
                                    </select>
                                </div>

                                <div class="form-group">
                                    <label>Clinic</label>
                                    <select id="medClinicSelect" required ${role === ROLES.EMPLOYEE ? 'disabled' : ''}>
                                        <option value="">Select clinic...</option>
                                        ${clinics.map(c => `
                                            <option value="${c.id}" ${c.id === myClinicId ? 'selected' : ''}>${this.escapeHtml(c.name)}</option>
                                        `).join('')}
                                    </select>
                                </div>
                            </div>

                            <div class="form-row">
                                <div class="form-group">
                                    <label>Diagnosis</label>
                                    <input type="text" id="diagnosis" required>
                                </div>

                                <div class="form-group">
                                    <label>Complaint</label>
                                    <input type="text" id="complaint" required>
                                </div>
                            </div>

                            <div class="form-row">
                                <div class="form-group">
                                    <label>Treatment</label>
                                    <textarea id="treatment" rows="2" required></textarea>
                                </div>

                                <div class="form-group">
                                    <label>Medications</label>
                                    <textarea id="medications" rows="2" required></textarea>
                                </div>
                            </div>

                            <div class="form-row">
                                <div class="form-group">
                                    <label>Allergies</label>
                                    <input type="text" id="allergies" value="Нет">
                                </div>

                                <div class="form-group">
                                    <label>Lab Results</label>
                                    <input type="text" id="labResults">
                                </div>
                            </div>

                            <div class="form-group full-width">
                                <label>Doctor Notes</label>
                                <textarea id="doctorNotes" rows="3"></textarea>
                            </div>

                            <div id="createMedicalError" style="color: red; margin: 10px 0;"></div>
                            <div id="createMedicalSuccess" style="color: green; margin: 10px 0;"></div>

                            <div class="form-actions">
                                <button type="submit" class="btn btn-primary">Create Medical Record</button>
                                <button type="button" class="btn btn-secondary" onclick="window.location.hash='medical-data'">Cancel</button>
                            </div>
                        </form>
                    </div>
                </main>
            </div>
        `;

        document.getElementById('createMedicalForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const patientSelect = document.getElementById('patientId');
            const clinicSelect = document.getElementById('medClinicSelect');

            const patientName = patientSelect.selectedOptions[0]?.textContent || 'Unknown';
            const clinicName = clinicSelect.selectedOptions[0]?.textContent || 'Unknown';

            const medicalData = {
                user_id: parseInt(patientSelect.value),
                clinic_id: parseInt(clinicSelect.value),
                diagnosis: document.getElementById('diagnosis').value,
                complaint: document.getElementById('complaint').value,
                treatment: document.getElementById('treatment').value,
                medications: document.getElementById('medications').value,
                allergies: document.getElementById('allergies').value,
                lab_results: document.getElementById('labResults').value,
                doctor_notes: document.getElementById('doctorNotes').value
            };

            try {
                // Оборачиваем в { medical_data: ... } — как ожидает бэкенд
                await api.createMedicalData({ medical_data: medicalData });

                this.logEvent('success', `Medical record created for ${patientName} at ${clinicName}`);
                document.getElementById('createMedicalSuccess').textContent = 'Medical record created successfully!';
                document.getElementById('createMedicalError').textContent = '';
                document.getElementById('createMedicalForm').reset();
            } catch (error) {
                this.logEvent('error', `Failed to create medical record: ${error.message}`);
                document.getElementById('createMedicalError').textContent = error.message || 'Failed to create record';
                document.getElementById('createMedicalSuccess').textContent = '';
            }
        });;
    }

    async renderCreateTransfer() {
        const role = auth.getRole();
        const myClinicId = auth.getClinicId();

        let medicalRecords = [];
        let clinics = [];

        try {
            const filters = role === ROLES.EMPLOYEE ? { clinic_id: myClinicId } : {};
            const response = await api.getMedicalData(filters);
            medicalRecords = response.medical_data || [];

            // Предзагружаем имена
            const userIds = [...new Set(medicalRecords.map(r => r.user_id))];
            await Promise.all(userIds.map(id => this.getUserName(id)));

            // Load clinics
            const clinicsResponse = await api.getClinics();
            clinics = clinicsResponse.clinics || [];

        } catch (error) {
            console.error('Failed to load data:', error);
        }

        const clinicOptions = clinics.map(clinic =>
            `<option value="${clinic.id}">${this.escapeHtml(clinic.name)}</option>`
        ).join('');

        this.appContainer.innerHTML = `
        <div class="dashboard-wrapper">
            ${this.renderSidebar()}

            <main class="main-content">
                <h1 style="margin-bottom: 30px;">Create Data Transfer Request</h1>

                <div class="widget" style="max-width: 600px;">
                    <form id="createTransferForm">
                        <div class="form-group">
                            <label>Medical Record to Transfer</label>
                            <select id="medicalDataId" required>
                                <option value="">Select medical record...</option>
                                ${medicalRecords.map(record => `
                                    <option value="${record.id}" data-user-id="${record.user_id}" data-clinic-id="${record.clinic_id}">
                                        ${this.escapeHtml(this.userNameCache.get(record.user_id) || 'Loading...')} - ${this.escapeHtml(record.diagnosis)} (${new Date(record.created_at).toLocaleDateString()})
                                    </option>
                                `).join('')}
                            </select>
                        </div>

                        <div class="form-group">
                            <label>From Clinic</label>
                            <select id="fromClinicSelect" required ${role === ROLES.EMPLOYEE ? 'disabled' : ''}>
                                <option value="">Select clinic...</option>
                                ${clinicOptions}
                            </select>
                        </div>

                        <div class="form-group">
                            <label>To Clinic</label>
                            <select id="toClinicSelect" required>
                                <option value="">Select target clinic...</option>
                                ${clinicOptions}
                            </select>
                        </div>

                        <div class="form-group">
                            <label>Patient (auto-filled)</label>
                            <input type="text" id="transferUserDisplay" readonly placeholder="Select medical record first">
                            <input type="hidden" id="transferUserId" required>
                        </div>

                        <div id="createTransferError" style="color: red; margin: 10px 0;"></div>
                        <div id="createTransferSuccess" style="color: green; margin: 10px 0;"></div>

                        <div class="form-actions">
                            <button type="submit" class="btn btn-primary">Create Transfer Request</button>
                            <button type="button" class="btn btn-secondary" onclick="window.location.hash='transfers'">Cancel</button>
                        </div>
                    </form>
                </div>
            </main>
        </div>
    `;

        // Pre-select from clinic for employees
        if (role === ROLES.EMPLOYEE) {
            document.getElementById('fromClinicSelect').value = myClinicId;
        }

        // Автозаполнение
        document.getElementById('medicalDataId').addEventListener('change', (e) => {
            const selectedOption = e.target.selectedOptions[0];
            if (selectedOption && selectedOption.value) {
                const userId = parseInt(selectedOption.dataset.userId);
                const clinicId = parseInt(selectedOption.dataset.clinicId);

                // Auto-fill patient
                const userName = this.userNameCache.get(userId) || 'Unknown';
                document.getElementById('transferUserDisplay').value = userName;
                document.getElementById('transferUserId').value = userId;

                // Auto-fill from clinic
                document.getElementById('fromClinicSelect').value = clinicId;
            }
        });

        document.getElementById('createTransferForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const fromClinicSelect = document.getElementById('fromClinicSelect');
            const toClinicSelect = document.getElementById('toClinicSelect');

            const fromClinicName = fromClinicSelect.selectedOptions[0]?.textContent || 'Unknown';
            const toClinicName = toClinicSelect.selectedOptions[0]?.textContent || 'Unknown';

            const transferData = {
                medical_data_id: parseInt(document.getElementById('medicalDataId').value),
                from_clinic_id: parseInt(fromClinicSelect.value),
                to_clinic_id: parseInt(toClinicSelect.value),
                user_id: parseInt(document.getElementById('transferUserId').value)
            };

            try {
                await api.createTransfer(transferData);
                this.logEvent('success', `Transfer request created from ${fromClinicName} to ${toClinicName}`);
                document.getElementById('createTransferSuccess').textContent = 'Transfer request created! Patient must approve it.';
                document.getElementById('createTransferError').textContent = '';
                document.getElementById('createTransferForm').reset();
            } catch (error) {
                this.logEvent('error', `Failed to create transfer: ${error.message}`);
                document.getElementById('createTransferError').textContent = error.message || 'Failed to create transfer';
                document.getElementById('createTransferSuccess').textContent = '';
            }
        });
    }

    logout() {
        auth.clearAuth();
        window.location.hash = 'login';
    }

    // Получить имя пользователя по ID (с кэшем)
    async getUserName(userId) {
        if (!userId) return 'Unknown';
        if (this.userNameCache.has(userId)) {
            return this.userNameCache.get(userId);
        }

        try {
            const response = await api.getUser(userId);
            const fullName = `${response.user.first_name} ${response.user.last_name}`;
            this.userNameCache.set(userId, fullName);
            return fullName;
        } catch (error) {
            console.warn(`Failed to fetch user ${userId}:`, error);
            const fallback = `User #${userId}`;
            this.userNameCache.set(userId, fallback);
            return fallback;
        }
    }

    // Безопасный вывод HTML
    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    async renderClinics() {
        const role = auth.getRole();
        let clinics = [];

        try {
            const response = await api.getClinics();
            clinics = response.clinics || [];
        } catch (error) {
            console.error('Failed to load clinics:', error);
        }

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px;">
                        <h1 style="margin: 0;">Clinics Management</h1>
                        ${role === ROLES.ADMIN ? '<button class="btn btn-primary" onclick="app.showClinicModal()"><i class="fas fa-plus"></i> Create Clinic</button>' : ''}
                    </div>

                    <div class="widget">
                        <table class="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Address</th>
                                    <th>Phone</th>
                                    <th>Email</th>
                                    <th>Created</th>
                                    ${role === ROLES.ADMIN ? '<th>Actions</th>' : ''}
                                </tr>
                            </thead>
                            <tbody>
                                ${clinics.length === 0 ? `
                                    <tr>
                                        <td colspan="${role === ROLES.ADMIN ? '7' : '6'}" style="text-align: center; padding: 20px; color: var(--text-secondary);">
                                            No clinics found
                                        </td>
                                    </tr>
                                ` : clinics.map(clinic => `
                                    <tr>
                                        <td>${clinic.id}</td>
                                        <td>${this.escapeHtml(clinic.name || 'N/A')}</td>
                                        <td>${this.escapeHtml(clinic.address || 'N/A')}</td>
                                        <td>${this.escapeHtml(clinic.phone || 'N/A')}</td>
                                        <td>${this.escapeHtml(clinic.email || 'N/A')}</td>
                                        <td>${new Date(clinic.created_at).toLocaleDateString()}</td>
                                        ${role === ROLES.ADMIN ? `
                                            <td>
                                                <button class="btn btn-primary btn-sm" onclick="app.showClinicModal(${clinic.id})"><i class="fas fa-edit"></i> Edit</button>
                                                <button class="btn btn-danger btn-sm" onclick="app.deleteClinic(${clinic.id})"><i class="fas fa-trash"></i> Delete</button>
                                            </td>
                                        ` : ''}
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    </div>
                </main>
            </div>

            <!-- Clinic Modal -->
            <div id="clinicModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); z-index: 2000; align-items: center; justify-content: center;">
                <div style="background: var(--bg-primary); padding: 32px; border-radius: 12px; max-width: 600px; width: 90%; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 50px var(--shadow);">
                    <h2 id="clinicModalTitle" style="margin-bottom: 24px; color: var(--text-primary);">Create Clinic</h2>
                    <form id="clinicForm">
                        <input type="hidden" id="clinicId">
                        <div class="form-group">
                            <label>Name</label>
                            <input type="text" id="clinicName" required>
                        </div>
                        <div class="form-group">
                            <label>Address</label>
                            <input type="text" id="clinicAddress" required>
                        </div>
                        <div class="form-group">
                            <label>Phone</label>
                            <input type="tel" id="clinicPhone" required>
                        </div>
                        <div class="form-group">
                            <label>Email</label>
                            <input type="email" id="clinicEmail" required>
                        </div>
                        <div class="form-actions">
                            <button type="submit" class="btn btn-primary"><i class="fas fa-save"></i> Save</button>
                            <button type="button" class="btn btn-secondary" onclick="app.hideClinicModal()"><i class="fas fa-times"></i> Cancel</button>
                        </div>
                    </form>
                </div>
            </div>
        `;

        // Setup form handler (remove old listeners by cloning)
        if (role === ROLES.ADMIN) {
            const form = document.getElementById('clinicForm');
            const newForm = form.cloneNode(true);
            form.parentNode.replaceChild(newForm, form);
            newForm.addEventListener('submit', (e) => this.submitClinicForm(e));
        }
    }

    async renderEventLog() {
        const logs = this.getFilteredEventLogs();
        const role = auth.getRole();

        this.appContainer.innerHTML = `
            <div class="dashboard-wrapper">
                ${this.renderSidebar()}

                <main class="main-content">
                    <h1 style="margin-bottom: 30px;">Event Log ${role === ROLES.ADMIN ? '(All Events)' : '(My Events)'}</h1>

                    <div class="widget">
                        <div style="margin-bottom: 15px;">
                            <button class="btn btn-secondary btn-sm" onclick="app.clearEventLogs()">Clear My Logs</button>
                        </div>
                        <div class="log-container" style="max-height: 600px; overflow-y: auto; font-family: monospace; font-size: 12px; background: var(--bg-secondary); padding: 15px; border-radius: 4px;">
                            ${logs.length === 0 ? `
                                <p style="text-align: center; color: var(--text-secondary);">No events logged yet</p>
                            ` : logs.map(log => `
                                <div style="margin-bottom: 10px; padding: 8px; background: var(--bg-primary); border-left: 3px solid ${this.getLogColor(log.level)}; border-radius: 4px;">
                                    <strong>[${new Date(log.timestamp).toLocaleString()}]</strong>
                                    <span style="color: ${this.getLogColor(log.level)}; font-weight: bold;">[${log.level.toUpperCase()}]</span>
                                    ${this.escapeHtml(log.message)}
                                </div>
                            `).join('')}
                        </div>
                    </div>
                </main>
            </div>
        `;
    }

    logEvent(level, message) {
        const user = auth.getUser();
        const role = auth.getRole();

        const logs = this.getAllEventLogs();
        logs.push({
            timestamp: new Date().toISOString(),
            level: level,
            message: message,
            user_id: user.id,
            user_email: user.email,
            clinic_id: role === ROLES.EMPLOYEE ? auth.getClinicId() : null
        });

        // Keep only last 1000 logs globally
        if (logs.length > 1000) {
            logs.shift();
        }
        localStorage.setItem('event_logs', JSON.stringify(logs));
    }

    getAllEventLogs() {
        const logs = localStorage.getItem('event_logs');
        return logs ? JSON.parse(logs) : [];
    }

    getFilteredEventLogs() {
        const allLogs = this.getAllEventLogs();
        const role = auth.getRole();
        const user = auth.getUser();

        if (role === ROLES.ADMIN) {
            // Admin sees all logs
            return allLogs;
        } else if (role === ROLES.EMPLOYEE) {
            // Employee sees only logs from their clinic
            const clinicId = auth.getClinicId();
            return allLogs.filter(log => log.clinic_id === clinicId || log.user_id === user.id);
        } else {
            // Patient sees only their own logs
            return allLogs.filter(log => log.user_id === user.id);
        }
    }

    getEventLogs() {
        // Alias for compatibility
        return this.getFilteredEventLogs();
    }

    clearEventLogs() {
        const role = auth.getRole();

        if (role === ROLES.ADMIN) {
            // Admin can clear all logs
            if (confirm('Are you sure you want to clear ALL event logs?')) {
                localStorage.removeItem('event_logs');
            }
        } else {
            // Other users clear only their logs
            const user = auth.getUser();
            const allLogs = this.getAllEventLogs();
            const clinicId = role === ROLES.EMPLOYEE ? auth.getClinicId() : null;

            const remainingLogs = allLogs.filter(log => {
                if (role === ROLES.EMPLOYEE) {
                    // Keep logs that are NOT from this employee or their clinic
                    return log.clinic_id !== clinicId && log.user_id !== user.id;
                } else {
                    // Keep logs that are NOT from this user
                    return log.user_id !== user.id;
                }
            });

            localStorage.setItem('event_logs', JSON.stringify(remainingLogs));
        }

        this.renderEventLog();
    }

    getLogColor(level) {
        switch(level.toLowerCase()) {
            case 'error': return '#dc3545';
            case 'warning': return '#ffc107';
            case 'info': return '#17a2b8';
            case 'success': return '#28a745';
            default: return '#6c757d';
        }
    }

    async showClinicModal(clinicId = null) {
        const modal = document.getElementById('clinicModal');
        const title = document.getElementById('clinicModalTitle');
        const form = document.getElementById('clinicForm');

        if (clinicId) {
            // Edit mode
            title.textContent = 'Edit Clinic';
            try {
                const response = await api.getClinic(clinicId);
                const clinic = response.clinic;
                document.getElementById('clinicId').value = clinic.id;
                document.getElementById('clinicName').value = clinic.name;
                document.getElementById('clinicAddress').value = clinic.address;
                document.getElementById('clinicPhone').value = clinic.phone;
                document.getElementById('clinicEmail').value = clinic.email;
            } catch (error) {
                alert('Failed to load clinic data: ' + error.message);
                return;
            }
        } else {
            // Create mode
            title.textContent = 'Create Clinic';
            form.reset();
            document.getElementById('clinicId').value = '';
        }

        modal.style.display = 'flex';
    }

    hideClinicModal() {
        const modal = document.getElementById('clinicModal');
        modal.style.display = 'none';
    }

    async submitClinicForm(e) {
        e.preventDefault();

        // Prevent double submission
        const submitButton = e.target.querySelector('button[type="submit"]');
        if (submitButton.disabled) {
            return;
        }
        submitButton.disabled = true;

        const clinicId = document.getElementById('clinicId').value;
        const clinicData = {
            name: document.getElementById('clinicName').value,
            address: document.getElementById('clinicAddress').value,
            phone: document.getElementById('clinicPhone').value,
            email: document.getElementById('clinicEmail').value
        };

        try {
            if (clinicId) {
                // Update existing clinic
                await api.updateClinic(clinicId, clinicData);
                this.logEvent('success', `Clinic "${clinicData.name}" updated successfully`);
                alert('Clinic updated successfully!');
            } else {
                // Create new clinic
                await api.createClinic(clinicData);
                this.logEvent('success', `Clinic "${clinicData.name}" created successfully`);
                alert('Clinic created successfully!');
            }

            this.hideClinicModal();
            this.handleRoute();
        } catch (error) {
            this.logEvent('error', `Failed to save clinic: ${error.message}`);
            alert('Failed to save clinic: ' + error.message);
            submitButton.disabled = false; // Re-enable on error
        }
    }

    async deleteClinic(clinicId) {
        if (!confirm('Are you sure you want to delete this clinic?')) {
            return;
        }

        try {
            // Try to get clinic name before deleting
            let clinicName = `Clinic #${clinicId}`;
            try {
                const response = await api.getClinic(clinicId);
                clinicName = response.clinic?.name || clinicName;
            } catch (e) {
                // If can't get name, use ID
            }

            await api.deleteClinic(clinicId);
            this.logEvent('success', `Clinic "${clinicName}" deleted successfully`);
            alert('Clinic deleted successfully!');
            this.handleRoute();
        } catch (error) {
            this.logEvent('error', `Failed to delete clinic #${clinicId}: ${error.message}`);
            alert('Failed to delete clinic: ' + error.message);
        }
    }

    async showUserModal(userType, userId = null) {
        const modal = document.getElementById('userModal');
        const title = document.getElementById('userModalTitle');
        const form = document.getElementById('userForm');
        const clinicGroup = document.getElementById('userClinicGroup');
        const passwordGroup = document.getElementById('userPasswordGroup');

        document.getElementById('userType').value = userType;

        if (userId) {
            // Edit mode
            title.textContent = `Edit ${userType.charAt(0).toUpperCase() + userType.slice(1)}`;
            passwordGroup.querySelector('input').required = false;

            try {
                const response = await api.getUser(userId);
                const user = response.user;
                document.getElementById('userId').value = user.id;
                document.getElementById('userFirstName').value = user.first_name;
                document.getElementById('userLastName').value = user.last_name;
                document.getElementById('userEmail').value = user.email;
                document.getElementById('userPhone').value = user.phone_number;
                document.getElementById('userPassword').value = '';

                if (userType === 'employee') {
                    clinicGroup.style.display = 'block';
                    document.getElementById('userClinicId').value = user.clinic_id?.value || user.clinic_id || '';
                } else {
                    clinicGroup.style.display = 'none';
                }
            } catch (error) {
                alert('Failed to load user data: ' + error.message);
                return;
            }
        } else {
            // Create mode
            title.textContent = `Create ${userType.charAt(0).toUpperCase() + userType.slice(1)}`;
            form.reset();
            document.getElementById('userId').value = '';
            passwordGroup.querySelector('input').required = true;

            const currentRole = auth.getRole();

            if (userType === 'employee') {
                clinicGroup.style.display = 'block';
                document.getElementById('userClinicId').required = true;

                // Auto-fill clinic_id for employee creating another employee
                if (currentRole === ROLES.EMPLOYEE) {
                    document.getElementById('userClinicId').value = auth.getClinicId();
                    document.getElementById('userClinicId').readOnly = true;
                }
            } else if (userType === 'patient') {
                // For patients, employee should also specify clinic_id
                if (currentRole === ROLES.EMPLOYEE) {
                    clinicGroup.style.display = 'block';
                    document.getElementById('userClinicId').required = true;
                    document.getElementById('userClinicId').value = auth.getClinicId();
                    document.getElementById('userClinicId').readOnly = true;
                } else {
                    clinicGroup.style.display = 'none';
                    document.getElementById('userClinicId').required = false;
                }
            } else {
                clinicGroup.style.display = 'none';
                document.getElementById('userClinicId').required = false;
            }
        }

        modal.style.display = 'flex';
    }

    hideUserModal() {
        const modal = document.getElementById('userModal');
        modal.style.display = 'none';
    }

    async submitUserForm(e) {
        e.preventDefault();

        const userId = document.getElementById('userId').value;
        const userType = document.getElementById('userType').value;
        const password = document.getElementById('userPassword').value;

        const userData = {
            first_name: document.getElementById('userFirstName').value,
            last_name: document.getElementById('userLastName').value,
            email: document.getElementById('userEmail').value,
            phone_number: document.getElementById('userPhone').value,
            type: userType
        };

        // Add password if provided (for create or update)
        if (password) {
            userData.password = password;
        }

        // Add clinic_id for employees AND patients when created by employee
        const clinicIdInput = document.getElementById('userClinicId');
        if (clinicIdInput.value) {
            userData.clinic_id = parseInt(clinicIdInput.value);
        }

        try {
            if (userId) {
                // Update existing user
                await api.updateUser(userId, userData);
                this.logEvent('success', `${userType} "${userData.first_name} ${userData.last_name}" updated successfully`);
                alert('User updated successfully!');
            } else {
                // Create new user
                await api.createUser(userData);
                this.logEvent('success', `${userType} "${userData.first_name} ${userData.last_name}" created successfully`);
                alert('User created successfully!');
            }

            this.hideUserModal();
            this.handleRoute();
        } catch (error) {
            this.logEvent('error', `Failed to save user: ${error.message}`);
            alert('Failed to save user: ' + error.message);
        }
    }

    async deleteUser(userId) {
        if (!confirm('Are you sure you want to delete this user?')) {
            return;
        }

        try {
            // Get user name before deleting
            const userName = this.userNameCache.get(userId) || await this.getUserName(userId);
            await api.deleteUser(userId);
            this.logEvent('success', `User "${userName}" deleted successfully`);
            alert('User deleted successfully!');
            // Refresh the current page
            this.handleRoute();
        } catch (error) {
            const userName = this.userNameCache.get(userId) || `User #${userId}`;
            this.logEvent('error', `Failed to delete "${userName}": ${error.message}`);
            alert('Failed to delete user: ' + error.message);
        }
    }

    async showMedicalDataModal() {
        const modal = document.getElementById('medicalDataModal');
        const form = document.getElementById('medicalDataForm');
        form.reset();

        const role = auth.getRole();
        const clinicSelect = document.getElementById('medicalDataClinicSelect');

        try {
            // Load clinics list
            const clinicsResponse = await api.getClinics();
            const clinics = clinicsResponse.clinics || [];

            clinicSelect.innerHTML = '<option value="">Select clinic...</option>';
            clinics.forEach(clinic => {
                const option = document.createElement('option');
                option.value = clinic.id;
                option.textContent = clinic.name;
                clinicSelect.appendChild(option);
            });

            // Pre-select and disable clinic for employees
            if (role === ROLES.EMPLOYEE) {
                const userClinicId = auth.getClinicId();
                clinicSelect.value = userClinicId;
                clinicSelect.disabled = true;
            } else {
                clinicSelect.disabled = false;
            }

            // Load patients list
            const response = await api.getUsers();
            const users = response.users || [];
            const patients = users.filter(u => u.type === ROLES.PATIENT);

            const patientSelect = document.getElementById('medicalDataPatientSelect');
            patientSelect.innerHTML = '<option value="">Select patient...</option>';

            patients.forEach(patient => {
                const option = document.createElement('option');
                option.value = patient.id;
                option.textContent = `${patient.first_name} ${patient.last_name} (${patient.email})`;
                patientSelect.appendChild(option);
            });
        } catch (error) {
            console.error('Failed to load data:', error);
            alert('Failed to load form data: ' + error.message);
        }

        modal.style.display = 'flex';
    }

    hideMedicalDataModal() {
        const modal = document.getElementById('medicalDataModal');
        modal.style.display = 'none';
    }

    async submitMedicalDataForm(e) {
        e.preventDefault();

        const patientSelect = document.getElementById('medicalDataPatientSelect');
        const clinicSelect = document.getElementById('medicalDataClinicSelect');

        const patientId = patientSelect.value;
        const patientName = patientSelect.selectedOptions[0].textContent;
        const clinicName = clinicSelect.selectedOptions[0].textContent;

        const medicalData = {
            user_id: parseInt(patientId),
            clinic_id: parseInt(clinicSelect.value),
            diagnosis: document.getElementById('medicalDataDiagnosis').value,
            complaint: document.getElementById('medicalDataComplaint').value,
            treatment: document.getElementById('medicalDataTreatment').value,
            medications: document.getElementById('medicalDataMedications').value
        };

        try {
            await api.createMedicalData({ medical_data: medicalData });
            this.logEvent('success', `Medical record created for ${patientName} at ${clinicName}`);
            alert('Medical record created successfully!');
            this.hideMedicalDataModal();
            this.handleRoute();
        } catch (error) {
            this.logEvent('error', `Failed to create medical record: ${error.message}`);
            alert('Failed to create medical record: ' + error.message);
        }
    }

    async deleteMedicalData(dataId) {
        if (!confirm('Are you sure you want to delete this medical record?')) {
            return;
        }

        try {
            await api.deleteMedicalData(dataId);
            this.logEvent('success', `Medical record #${dataId} deleted successfully`);
            alert('Medical record deleted successfully!');
            this.handleRoute();
        } catch (error) {
            this.logEvent('error', `Failed to delete medical record #${dataId}: ${error.message}`);
            alert('Failed to delete medical record: ' + error.message);
        }
    }

    logout() {
        this.userNameCache.clear(); // Очистка кэша при выходе
        auth.clearAuth();
        window.location.hash = 'login';
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.app = new App();
});
