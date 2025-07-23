<<<<<<< HEAD
class WebScraperApp {
    constructor() {
        this.apiBase = '/api';
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadResults();
        this.checkHealth();
        setInterval(() => this.checkHealth(), 30000);
    }

    bindEvents() {
        const events = {
            'scrapeForm': ['submit', e => { e.preventDefault(); this.handleScrape(); }],
            'refreshBtn': ['click', () => this.loadResults()],
            'closeModal': ['click', () => this.closeModal()],
            'detailModal': ['click', e => e.target.id === 'detailModal' && this.closeModal()]
        };

        Object.entries(events).forEach(([id, [event, handler]]) => 
            document.getElementById(id)?.addEventListener(event, handler)
        );

        document.addEventListener('keydown', e => e.key === 'Escape' && this.closeModal());
    }

    async handleScrape() {
        const urlInput = document.getElementById('urlInput');
        const scrapeBtn = document.getElementById('scrapeBtn');
        const loadingIndicator = document.getElementById('loadingIndicator');
        const url = urlInput.value.trim();
        
        if (!url || !this.isValidUrl(url)) {
            this.showAlert('Por favor, ingresa una URL válida con http:// o https://', 'error');
            return;
        }

        this.toggleLoading(scrapeBtn, loadingIndicator, true);

        try {
            const response = await fetch(`${this.apiBase}/scrape`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url })
            });
            const data = await response.json();

            if (response.ok) {
                this.showAlert(`Se ha scrapeado el sitio web con éxito: ${url}`, 'success');
                urlInput.value = '';
                this.loadResults();
            } else {
                this.showAlert(data.error || 'Fallo al scrapear la URL', 'error');
            }
        } catch (error) {
            console.error('Scraping error:', error);
            this.showAlert('Ha ocurrido un error de red. Por favor, inténtalo de nuevo.', 'error');
        } finally {
            this.toggleLoading(scrapeBtn, loadingIndicator, false);
        }
    }

    toggleLoading(btn, indicator, isLoading) {
        btn.disabled = isLoading;
        btn.innerHTML = isLoading 
            ? `<span class="flex items-center justify-center space-x-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                <span>Scraping...</span>
               </span>`
            : `<span class="flex items-center justify-center space-x-2">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
                <span>Scrapear sitio web</span>
               </span>`;
        indicator.classList.toggle('hidden', !isLoading);
    }

    async apiRequest(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.apiBase}${endpoint}`, options);
            const data = await response.json();
            return { response, data };
        } catch (error) {
            console.error(`API request error (${endpoint}):`, error);
            throw error;
        }
    }

    async loadResults() {
        try {
            const { response, data } = await this.apiRequest('/results');
            if (response.ok) {
                this.displayResults(data.data || []);
            } else {
                this.showAlert('Fallo al cargar resultados', 'error');
            }
        } catch (error) {
            this.showAlert('Fallo al conectar con el servidor', 'error');
        }
    }

    displayResults(results) {
        const container = document.getElementById('resultsContainer');
        
        if (!results?.length) {
            container.innerHTML = this.getEmptyStateHTML();
            return;
        }
        container.innerHTML = results.map(result => this.createResultCard(result)).join('');
    }

    getEmptyStateHTML() {
        return `
            <div class="text-center text-gray-400 py-8">
                <svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                <p>Sin resultados aún. Comienza ingresando una URL arriba.</p>
            </div>
        `;
    }

    createResultCard(result) {
        const statusColor = result.status_code === 200 ? 'text-green-400' : 'text-red-400';
        const date = new Date(result.created_at).toLocaleString();
        const stats = [
            { icon: 'M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1', text: `${(result.links || []).length} links` },
            { icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z', text: `${(result.images || []).length} imágenes` },
            { icon: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253', text: `${result.word_count || 0} palabras` }
        ];
        
        return `
            <div class="bg-white/5 backdrop-blur-sm rounded-xl border border-white/20 p-6 hover:bg-white/10 transition-all duration-300 group">
                <div class="flex items-start justify-between mb-4">
                    <div class="flex-1 min-w-0">
                        <h3 class="text-lg font-semibold text-white truncate mb-2" title="${result.title || 'No title'}">
                            ${result.title || 'No title'}
                        </h3>
                        <p class="text-sm text-blue-400 truncate mb-2" title="${result.url}">${result.url}</p>
                        <p class="text-sm text-gray-400 mb-3 line-clamp-2">${result.description || 'No description available'}</p>
                    </div>
                    <span class="px-2 py-1 text-xs font-medium ${statusColor} bg-white/10 rounded-full ml-4">${result.status_code}</span>
                </div>
                
                <div class="flex items-center justify-between text-sm text-gray-400 mb-4">
                    <div class="flex items-center space-x-4">
                        ${stats.map(stat => `
                            <span class="flex items-center space-x-1">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${stat.icon}"></path>
                                </svg>
                                <span>${stat.text}</span>
                            </span>
                        `).join('')}
                    </div>
                    <span>${date}</span>
                </div>
                
                <div class="flex items-center justify-between">
                    <div class="flex items-center text-sm text-gray-500">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span>Tiempo de carga: ${result.load_time_ms}ms</span>
                    </div>
                    <div class="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button onclick="app.viewDetails(${result.id})" class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">Ver detalles</button>
                        <button onclick="app.deleteResult(${result.id})" class="px-3 py-1 bg-red-600 hover:bg-red-700 text-white text-sm rounded-lg transition-colors">Eliminar</button>
                    </div>
                </div>
            </div>
        `;
    }

    async viewDetails(id) {
        try {
            const { response, data } = await this.apiRequest(`/results/${id}`);
            if (response.ok) {
                this.showModal(data.data);
            } else {
                this.showAlert('Fallo al cargar detalles', 'error');
            }
        } catch (error) {
            this.showAlert('Fallo al cargar detalles', 'error');
        }
    }

    async deleteResult(id) {
        if (!confirm('¿Estás seguro de que quieres eliminar este resultado?')) return;

        try {
            const response = await fetch(`${this.apiBase}/results/${id}`, { method: 'DELETE' });
            
            if (response.ok) {
                this.showAlert('Resultado eliminado con éxito', 'success');
                this.loadResults();
            } else {
                let errorMessage = 'Fallo al eliminar resultado';
                const text = await response.text();
                
                if (text.trim()) {
                    const data = JSON.parse(text);
                    errorMessage = data.error || errorMessage;
                }
                this.showAlert(errorMessage, 'error');
            }
        } catch (error) {
            console.error('Delete error:', error);
            this.showAlert('Fallo al eliminar resultado', 'error');
        }
    }

    showModal(result) {
        const modal = document.getElementById('detailModal');
        const content = document.getElementById('modalContent');
        
        content.innerHTML = this.createDetailView(result);
        modal.classList.remove('hidden');
        document.body.style.overflow = 'hidden';
    }

    closeModal() {
        document.getElementById('detailModal').classList.add('hidden');
        document.body.style.overflow = 'auto';
    }

    createDetailView(result) {
        const date = new Date(result.created_at).toLocaleString();
        const statusColor = result.status_code === 200 ? 'text-green-400' : 'text-red-400';
        
        const sections = [
            {
                title: 'Información Básica',
                content: this.createBasicInfoSection(result, statusColor, date)
            },
            {
                title: 'Metadatos',
                content: this.createMetadataSection(result)
            },
            ...(result.headers?.length ? [{
                title: `Cabeceras (${result.headers.length})`,
                content: this.createListSection(result.headers, h => `<span class="text-purple-400 font-mono">H${h.level}:</span> <span class="text-white">${this.escapeHtml(h.text)}</span>`, 40)
            }] : []),
            ...(result.links?.length ? [{
                title: `Links (${result.links.length})`,
                content: this.createListSection(result.links, link => `<a href="${this.escapeHtml(link)}" target="_blank" rel="noopener noreferrer" class="text-blue-400 hover:text-blue-300 transition-colors truncate block">${this.escapeHtml(link)}</a>`, 50)
            }] : []),
            ...(result.images?.length ? [{
                title: `Imágenes (${result.images.length})`,
                content: this.createListSection(result.images, img => `<a href="${this.escapeHtml(img)}" target="_blank" rel="noopener noreferrer" class="text-green-400 hover:text-green-300 transition-colors truncate block">${this.escapeHtml(img)}</a>`, 20)
            }] : []),
            ...(result.content ? [{
                title: 'Content Preview',
                content: `<pre class="whitespace-pre-wrap break-words text-sm text-gray-300">${this.escapeHtml(result.content.substring(0, 2000))}${result.content.length > 2000 ? '...\n\n[Content truncated]' : ''}</pre>`
            }] : [])
        ];

        return `<div class="space-y-6">${sections.map(section => `
            <div class="bg-white/5 rounded-lg p-4">
                <h4 class="text-lg font-semibold text-white mb-3">${section.title}</h4>
                ${section.content}
            </div>
        `).join('')}</div>`;
    }

    createBasicInfoSection(result, statusColor, date) {
        const fields = [
            ['URL', result.url, 'text-blue-400 break-all'],
            ['Código de estado', result.status_code, `${statusColor} font-medium`],
            ['Tipo de contenido', result.content_type || 'Desconocido', 'text-white'],
            ['Tiempo de carga', `${result.load_time_ms}ms`, 'text-white'],
            ['Número de palabras', result.word_count || 0, 'text-white'],
            ['Scraped', date, 'text-white']
        ];

        return `<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            ${fields.map(([label, value, className]) => `
                <div>
                    <span class="text-gray-400">${label}:</span>
                    <p class="${className}">${value}</p>
                </div>
            `).join('')}
        </div>`;
    }

    createMetadataSection(result) {
        const fields = [
            ['Título', result.title || 'Sin título'],
            ['Descripción', result.description || 'Sin descripción'],
            ['Palabras clave', result.keywords || 'Sin palabras clave'],
            ['Autor', result.author || 'Desconocido'],
            ['Idioma', result.language || 'Desconocido'],
            ['Nombre del sitio', result.site_name || 'Desconocido']
        ];

        return `<div class="space-y-3 text-sm">
            ${fields.map(([label, value]) => `
                <div>
                    <span class="text-gray-400">${label}:</span>
                    <p class="text-white">${value}</p>
                </div>
            `).join('')}
        </div>`;
    }

    createListSection(items, formatter, limit) {
        const displayItems = items.slice(0, limit);
        const hasMore = items.length > limit;
        
        return `
            <div class="space-y-2 max-h-40 overflow-y-auto text-sm">
                ${displayItems.map(item => `<div>${formatter(item)}</div>`).join('')}
                ${hasMore ? `<div class="text-gray-500 text-xs mt-2">... and ${items.length - limit} more items</div>` : ''}
            </div>
        `;
    }

    async checkHealth() {
        try {
            const { response, data } = await this.apiRequest('/health');
            const isOnline = response.ok && data.data?.status === 'ok';
            this.updateStatus(isOnline);
        } catch (error) {
            this.updateStatus(false);
        }
    }

    updateStatus(isOnline) {
        const statusIndicator = document.getElementById('status-indicator');
        const statusText = statusIndicator.nextSibling;
        
        statusIndicator.className = `inline-block w-2 h-2 ${isOnline ? 'bg-green-400' : 'bg-red-400'} rounded-full animate-pulse mr-2`;
        statusText.textContent = ` Service ${isOnline ? 'Online' : 'Offline'}`;
    }

    showAlert(message, type = 'info') {
        const container = document.getElementById('alertContainer');
        const alertId = 'alert-' + Date.now();
        
        const config = {
            success: { class: 'bg-green-500/20 border-green-500/50 text-green-300', icon: 'M5 13l4 4L19 7' },
            error: { class: 'bg-red-500/20 border-red-500/50 text-red-300', icon: 'M6 18L18 6M6 6l12 12' },
            warning: { class: 'bg-yellow-500/20 border-yellow-500/50 text-yellow-300', icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 13.5c-.77.833.192 2.5 1.732 2.5z' },
            info: { class: 'bg-blue-500/20 border-blue-500/50 text-blue-300', icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' }
        }[type];

        const alert = document.createElement('div');
        alert.id = alertId;
        alert.className = `${config.class} border rounded-lg p-4 mb-4 backdrop-blur-sm flex items-center space-x-3 animate-in slide-in-from-top duration-300`;
        
        alert.innerHTML = `
            <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${config.icon}"></path>
            </svg>
            <span class="flex-1">${this.escapeHtml(message)}</span>
            <button onclick="document.getElementById('${alertId}').remove()" class="text-current hover:opacity-70 transition-opacity">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
            </button>
        `;

        container.appendChild(alert);
        setTimeout(() => document.getElementById(alertId)?.remove(), 5000);
    }

    isValidUrl(string) {
        try {
            const url = new URL(string);
            return ['http:', 'https:'].includes(url.protocol);
        } catch { return false; }
    }

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new WebScraperApp();
});

['error', 'unhandledrejection'].forEach(event => {
    window.addEventListener(event, e => {
        console.error(`Global ${event}:`, e.error || e.reason);
        app?.showAlert('Un error inesperado ha ocurrido', 'error');
    });
=======
class WebScraperApp {
    constructor() {
        this.apiBase = '/api';
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadResults();
        this.checkHealth();
        setInterval(() => this.checkHealth(), 30000);
    }

    bindEvents() {
        const events = {
            'scrapeForm': ['submit', e => { e.preventDefault(); this.handleScrape(); }],
            'refreshBtn': ['click', () => this.loadResults()],
            'closeModal': ['click', () => this.closeModal()],
            'detailModal': ['click', e => e.target.id === 'detailModal' && this.closeModal()],
        };

        Object.entries(events).forEach(([id, [event, handler]]) => 
            document.getElementById(id)?.addEventListener(event, handler)
        );

        document.addEventListener('keydown', e => e.key === 'Escape' && this.closeModal());
    }

    async handleScrape() {
        const urlInput = document.getElementById('urlInput');
        const scrapeBtn = document.getElementById('scrapeBtn');
        const loadingIndicator = document.getElementById('loadingIndicator');
        const url = urlInput.value.trim();
        
        if (!url || !this.isValidUrl(url)) {
            this.showAlert('Por favor, ingresa una URL válida con http:// o https://', 'error');
            return;
        }

        this.toggleLoading(scrapeBtn, loadingIndicator, true);

        try {
            const response = await fetch(`${this.apiBase}/scrape`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url })
            });
            const data = await response.json();

            if (response.ok) {
                this.showAlert(`Se ha scrapeado el sitio web con éxito: ${url}`, 'success');
                urlInput.value = '';
                this.loadResults();
            } else {
                this.showAlert(data.error || 'Fallo al scrapear la URL', 'error');
            }
        } catch (error) {
            console.error('Scraping error:', error);
            this.showAlert('Ha ocurrido un error de red. Por favor, inténtalo de nuevo.', 'error');
        } finally {
            this.toggleLoading(scrapeBtn, loadingIndicator, false);
        }
    }

    toggleLoading(btn, indicator, isLoading) {
        btn.disabled = isLoading;
        btn.innerHTML = isLoading 
            ? `<span class="flex items-center justify-center space-x-2">
                <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                <span>Scraping...</span>
               </span>`
            : `<span class="flex items-center justify-center space-x-2">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
                <span>Scrapear sitio web</span>
               </span>`;
        indicator.classList.toggle('hidden', !isLoading);
    }

    async apiRequest(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.apiBase}${endpoint}`, options);
            const data = await response.json();
            return { response, data };
        } catch (error) {
            console.error(`API request error (${endpoint}):`, error);
            throw error;
        }
    }

    async loadResults() {
        try {
            const { response, data } = await this.apiRequest('/results');
            if (response.ok) {
                this.displayResults(data.data || []);
            } else {
                this.showAlert('Fallo al cargar resultados', 'error');
            }
        } catch (error) {
            this.showAlert('Fallo al conectar con el servidor', 'error');
        }
    }

    displayResults(results) {
        const container = document.getElementById('resultsContainer');
        
        if (!results?.length) {
            container.innerHTML = this.getEmptyStateHTML();
            return;
        }
        container.innerHTML = results.map(result => this.createResultCard(result)).join('');
    }

    getEmptyStateHTML() {
        return `
            <div class="text-center text-gray-400 py-8">
                <svg class="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                <p>Sin resultados aún. Comienza ingresando una URL arriba.</p>
            </div>
        `;
    }

    createResultCard(result) {
        const statusColor = result.status_code === 200 ? 'text-green-400' : 'text-red-400';
        const date = new Date(result.created_at).toLocaleString();
        const stats = [
            { icon: 'M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1', text: `${(result.links || []).length} links` },
            { icon: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z', text: `${(result.images || []).length} imágenes` },
            { icon: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253', text: `${result.word_count || 0} palabras` }
        ];
        
        return `
            <div class="bg-white/5 backdrop-blur-sm rounded-xl border border-white/20 p-6 hover:bg-white/10 transition-all duration-300 group">
                <div class="flex items-start justify-between mb-4">
                    <div class="flex-1 min-w-0">
                        <h3 class="text-lg font-semibold text-white truncate mb-2" title="${result.title || 'No title'}">
                            ${result.title || 'No title'}
                        </h3>
                        <p class="text-sm text-blue-400 truncate mb-2" title="${result.url}">${result.url}</p>
                        <p class="text-sm text-gray-400 mb-3 line-clamp-2">${result.description || 'No description available'}</p>
                    </div>
                    <span class="px-2 py-1 text-xs font-medium ${statusColor} bg-white/10 rounded-full ml-4">${result.status_code}</span>
                </div>
                
                <div class="flex items-center justify-between text-sm text-gray-400 mb-4">
                    <div class="flex items-center space-x-4">
                        ${stats.map(stat => `
                            <span class="flex items-center space-x-1">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${stat.icon}"></path>
                                </svg>
                                <span>${stat.text}</span>
                            </span>
                        `).join('')}
                    </div>
                    <span>${date}</span>
                </div>
                
                <div class="flex items-center justify-between">
                    <div class="flex items-center text-sm text-gray-500">
                        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span>Tiempo de carga: ${result.load_time_ms}ms</span>
                    </div>
                    <div class="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button onclick="app.viewDetails(${result.id})" class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors">Ver detalles</button>
                        <button onclick="app.deleteResult(${result.id})" class="px-3 py-1 bg-red-600 hover:bg-red-700 text-white text-sm rounded-lg transition-colors">Eliminar</button>
                    </div>
                </div>
            </div>
        `;
    }

    async viewDetails(id) {
        try {
            const { response, data } = await this.apiRequest(`/results/${id}`);
            if (response.ok) {
                this.showModal(data.data);
            } else {
                this.showAlert('Fallo al cargar detalles', 'error');
            }
        } catch (error) {
            this.showAlert('Fallo al cargar detalles', 'error');
        }
    }

    async deleteResult(id) {
        if (!confirm('¿Estás seguro de que quieres eliminar este resultado?')) return;

        try {
            const response = await fetch(`${this.apiBase}/results/${id}`, { method: 'DELETE' });
            
            if (response.ok) {
                this.showAlert('Resultado eliminado con éxito', 'success');
                this.loadResults();
            } else {
                let errorMessage = 'Fallo al eliminar resultado';
                const text = await response.text();
                
                if (text.trim()) {
                    const data = JSON.parse(text);
                    errorMessage = data.error || errorMessage;
                }
                this.showAlert(errorMessage, 'error');
            }
        } catch (error) {
            console.error('Delete error:', error);
            this.showAlert('Fallo al eliminar resultado', 'error');
        }
    }

    showModal(result) {
        const modal = document.getElementById('detailModal');
        const content = document.getElementById('modalContent');
        
        content.innerHTML = this.createDetailView(result);
        modal.classList.remove('hidden');
        document.body.style.overflow = 'hidden';
    }

    closeModal() {
        document.getElementById('detailModal').classList.add('hidden');
        document.body.style.overflow = 'auto';
    }

    createDetailView(result) {
        const date = new Date(result.created_at).toLocaleString();
        const statusColor = result.status_code === 200 ? 'text-green-400' : 'text-red-400';
        
        const sections = [
            {
                title: 'Información Básica',
                content: this.createBasicInfoSection(result, statusColor, date)
            },
            {
                title: 'Metadatos',
                content: this.createMetadataSection(result)
            },
            ...(result.headers?.length ? [{
                title: `Cabeceras (${result.headers.length})`,
                content: this.createListSection(result.headers, h => `<span class="text-purple-400 font-mono">H${h.level}:</span> <span class="text-white">${this.escapeHtml(h.text)}</span>`, 40)
            }] : []),
            ...(result.links?.length ? [{
                title: `Links (${result.links.length})`,
                content: this.createListSection(result.links, link => `<a href="${this.escapeHtml(link)}" target="_blank" rel="noopener noreferrer" class="text-blue-400 hover:text-blue-300 transition-colors truncate block">${this.escapeHtml(link)}</a>`, 50)
            }] : []),
            ...(result.images?.length ? [{
                title: `Imágenes (${result.images.length})`,
                content: this.createListSection(result.images, img => `<a href="${this.escapeHtml(img)}" target="_blank" rel="noopener noreferrer" class="text-green-400 hover:text-green-300 transition-colors truncate block">${this.escapeHtml(img)}</a>`, 20)
            }] : []),
            ...(result.content ? [{
                title: 'Content Preview',
                content: `<pre class="whitespace-pre-wrap break-words text-sm text-gray-300">${this.escapeHtml(result.content.substring(0, 2000))}${result.content.length > 2000 ? '...\n\n[Content truncated]' : ''}</pre>`
            }] : [])
        ];

        return `<div class="space-y-6">${sections.map(section => `
            <div class="bg-white/5 rounded-lg p-4">
                <h4 class="text-lg font-semibold text-white mb-3">${section.title}</h4>
                ${section.content}
            </div>
        `).join('')}</div>`;
    }

    createBasicInfoSection(result, statusColor, date) {
        const fields = [
            ['URL', result.url, 'text-blue-400 break-all'],
            ['Código de estado', result.status_code, `${statusColor} font-medium`],
            ['Tipo de contenido', result.content_type || 'Desconocido', 'text-white'],
            ['Tiempo de carga', `${result.load_time_ms}ms`, 'text-white'],
            ['Número de palabras', result.word_count || 0, 'text-white'],
            ['Scraped', date, 'text-white']
        ];

        return `<div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            ${fields.map(([label, value, className]) => `
                <div>
                    <span class="text-gray-400">${label}:</span>
                    <p class="${className}">${value}</p>
                </div>
            `).join('')}
        </div>`;
    }

    createMetadataSection(result) {
        const fields = [
            ['Título', result.title || 'Sin título'],
            ['Descripción', result.description || 'Sin descripción'],
            ['Palabras clave', result.keywords || 'Sin palabras clave'],
            ['Autor', result.author || 'Desconocido'],
            ['Idioma', result.language || 'Desconocido'],
            ['Nombre del sitio', result.site_name || 'Desconocido']
        ];

        return `<div class="space-y-3 text-sm">
            ${fields.map(([label, value]) => `
                <div>
                    <span class="text-gray-400">${label}:</span>
                    <p class="text-white">${value}</p>
                </div>
            `).join('')}
        </div>`;
    }

    createListSection(items, formatter, limit) {
        const displayItems = items.slice(0, limit);
        const hasMore = items.length > limit;
        
        return `
            <div class="space-y-2 max-h-40 overflow-y-auto text-sm">
                ${displayItems.map(item => `<div>${formatter(item)}</div>`).join('')}
                ${hasMore ? `<div class="text-gray-500 text-xs mt-2">... and ${items.length - limit} more items</div>` : ''}
            </div>
        `;
    }

    async checkHealth() {
        try {
            const { response, data } = await this.apiRequest('/health');
            const isOnline = response.ok && data.data?.status === 'ok';
            this.updateStatus(isOnline);
        } catch (error) {
            this.updateStatus(false);
        }
    }

    updateStatus(isOnline) {
        const statusIndicator = document.getElementById('status-indicator');
        const statusText = statusIndicator.nextSibling;
        
        statusIndicator.className = `inline-block w-2 h-2 ${isOnline ? 'bg-green-400' : 'bg-red-400'} rounded-full animate-pulse mr-2`;
        statusText.textContent = ` Service ${isOnline ? 'Online' : 'Offline'}`;
    }

    showAlert(message, type = 'info') {
        const container = document.getElementById('alertContainer');
        const alertId = 'alert-' + Date.now();
        
        const config = {
            success: { class: 'bg-green-500/20 border-green-500/50 text-green-300', icon: 'M5 13l4 4L19 7' },
            error: { class: 'bg-red-500/20 border-red-500/50 text-red-300', icon: 'M6 18L18 6M6 6l12 12' },
            warning: { class: 'bg-yellow-500/20 border-yellow-500/50 text-yellow-300', icon: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 13.5c-.77.833.192 2.5 1.732 2.5z' },
            info: { class: 'bg-blue-500/20 border-blue-500/50 text-blue-300', icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z' }
        }[type];

        const alert = document.createElement('div');
        alert.id = alertId;
        alert.className = `${config.class} border rounded-lg p-4 mb-4 backdrop-blur-sm flex items-center space-x-3 animate-in slide-in-from-top duration-300`;
        
        alert.innerHTML = `
            <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="${config.icon}"></path>
            </svg>
            <span class="flex-1">${this.escapeHtml(message)}</span>
            <button onclick="document.getElementById('${alertId}').remove()" class="text-current hover:opacity-70 transition-opacity">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
            </button>
        `;

        container.appendChild(alert);
        setTimeout(() => document.getElementById(alertId)?.remove(), 5000);
    }

    isValidUrl(string) {
        try {
            const url = new URL(string);
            return ['http:', 'https:'].includes(url.protocol);
        } catch { return false; }
    }

    escapeHtml(text) {
        if (!text) return '';
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new WebScraperApp();
});

['error', 'unhandledrejection'].forEach(event => {
    window.addEventListener(event, e => {
        console.error(`Global ${event}:`, e.error || e.reason);
        app?.showAlert('Un error inesperado ha ocurrido', 'error');
    });
>>>>>>> master
});