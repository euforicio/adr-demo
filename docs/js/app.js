// Main ADR Application
window.ADRApp = (function() {
    'use strict';

    // Application state
    const state = {
        adrs: [],
        filteredAdrs: [],
        currentAdr: null,
        searchTerm: '',
        statusFilter: '',
        isLoading: false,
        isMobileMenuOpen: false
    };

    // DOM elements
    const elements = {
        loading: null,
        app: null,
        sidebar: null,
        sidebarToggle: null,
        searchInput: null,
        searchClear: null,
        statusFilter: null,
        adrList: null,
        mainContent: null,
        welcomeScreen: null,
        adrContent: null,
        adrDocument: null,
        errorScreen: null,
        mobileOverlay: null,
        breadcrumb: null,
        tocList: null,
        printBtn: null,
        fullscreenBtn: null
    };

    // Statistics
    const stats = {
        totalAdrs: 0,
        acceptedAdrs: 0,
        c4Diagrams: 0
    };

    /**
     * Initialize the application
     */
    function init() {
        try {
            initializeElements();
            bindEvents();
            loadADRs();
            
            // Handle initial URL hash
            const hash = window.location.hash;
            if (hash && hash.startsWith('#adr-')) {
                const adrNumber = hash.replace('#adr-', '');
                setTimeout(() => loadADR(adrNumber), 100);
            }
        } catch (error) {
            console.error('Failed to initialize ADR App:', error);
            showError('Failed to initialize application');
        }
    }

    /**
     * Initialize DOM element references
     */
    function initializeElements() {
        const elementIds = [
            'loading', 'app', 'sidebar', 'sidebar-toggle', 'search-input', 
            'search-clear', 'status-filter', 'adr-list', 'main-content',
            'welcome-screen', 'adr-content', 'adr-document', 'error-screen',
            'mobile-overlay', 'current-adr-breadcrumb', 'toc-list', 'print-btn',
            'fullscreen-btn'
        ];

        elementIds.forEach(id => {
            const elementKey = id.replace(/-([a-z])/g, (g) => g[1].toUpperCase());
            elements[elementKey] = document.getElementById(id);
        });

        // Verify critical elements exist
        const criticalElements = ['loading', 'app', 'sidebar', 'adrList', 'mainContent'];
        criticalElements.forEach(key => {
            if (!elements[key]) {
                throw new Error(`Critical element not found: ${key}`);
            }
        });
    }

    /**
     * Bind event listeners
     */
    function bindEvents() {
        // Sidebar toggle
        if (elements.sidebarToggle) {
            elements.sidebarToggle.addEventListener('click', toggleMobileMenu);
        }

        // Mobile overlay
        if (elements.mobileOverlay) {
            elements.mobileOverlay.addEventListener('click', closeMobileMenu);
        }

        // Search functionality
        if (elements.searchInput) {
            elements.searchInput.addEventListener('input', handleSearch);
            elements.searchInput.addEventListener('keydown', handleSearchKeydown);
        }

        if (elements.searchClear) {
            elements.searchClear.addEventListener('click', clearSearch);
        }

        // Status filter
        if (elements.statusFilter) {
            elements.statusFilter.addEventListener('change', handleStatusFilter);
        }

        // Print button
        if (elements.printBtn) {
            elements.printBtn.addEventListener('click', printCurrentADR);
        }

        // Fullscreen button
        if (elements.fullscreenBtn) {
            elements.fullscreenBtn.addEventListener('click', toggleFullscreen);
        }

        // Window events
        window.addEventListener('hashchange', handleHashChange);
        window.addEventListener('resize', handleResize);
        
        // Keyboard shortcuts
        document.addEventListener('keydown', handleKeyboardShortcuts);

        // Handle browser back/forward
        window.addEventListener('popstate', handleHashChange);

        // Home navigation
        const homeLink = document.getElementById('home-link');
        const homeBreadcrumb = document.querySelector('.home-breadcrumb');
        
        if (homeLink) {
            homeLink.addEventListener('click', goHome);
            homeLink.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    goHome();
                }
            });
        }
        
        if (homeBreadcrumb) {
            homeBreadcrumb.addEventListener('click', goHome);
            homeBreadcrumb.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    goHome();
                }
            });
        }

        // Documentation buttons
        const showDocsBtn = document.getElementById('show-docs-btn');
        const showTemplateBtn = document.getElementById('show-template-btn');
        const docsLink = document.getElementById('docs-link');
        
        if (showDocsBtn) {
            showDocsBtn.addEventListener('click', showDocumentation);
        }
        
        if (showTemplateBtn) {
            showTemplateBtn.addEventListener('click', showTemplate);
        }
        
        if (docsLink) {
            docsLink.addEventListener('click', showDocumentation);
        }
    }

    /**
     * Load ADRs from the file system
     */
    async function loadADRs() {
        setLoading(true);
        
        try {
            // Get ADR list from README or scan directory
            const adrData = await ADRLoader.loadADRList();
            
            state.adrs = adrData;
            state.filteredAdrs = [...adrData];
            
            updateStatistics();
            renderADRList();
            showWelcomeScreen();
            
        } catch (error) {
            console.error('Failed to load ADRs:', error);
            showError('Failed to load Architecture Decision Records');
        } finally {
            setLoading(false);
        }
    }

    /**
     * Load and display a specific ADR
     */
    async function loadADR(adrNumber) {
        const adr = state.adrs.find(a => a.number === adrNumber);
        if (!adr) {
            showError(`ADR ${adrNumber} not found`);
            return;
        }

        setLoading(true);
        
        try {
            const content = await ADRLoader.loadADRContent(adr);
            state.currentAdr = { ...adr, content };
            
            renderADRContent(state.currentAdr);
            updateActiveNavigation(adrNumber);
            updateBreadcrumb(adr);
            
            // Update URL hash
            if (window.location.hash !== `#adr-${adrNumber}`) {
                window.history.pushState(null, null, `#adr-${adrNumber}`);
            }
            
        } catch (error) {
            console.error(`Failed to load ADR ${adrNumber}:`, error);
            showError(`Failed to load ADR ${adrNumber}`);
        } finally {
            setLoading(false);
        }
    }

    /**
     * Render the ADR list in the sidebar
     */
    function renderADRList() {
        if (!elements.adrList) return;

        elements.adrList.innerHTML = '';

        if (state.filteredAdrs.length === 0) {
            const noResults = document.createElement('li');
            noResults.className = 'no-results';
            noResults.innerHTML = `
                <div style="padding: var(--spacing-lg); text-align: center; color: var(--text-muted);">
                    <p>No ADRs found</p>
                    ${state.searchTerm || state.statusFilter ? 
                        '<button onclick="ADRApp.clearFilters()" style="margin-top: var(--spacing-sm); padding: var(--spacing-xs) var(--spacing-sm); border: 1px solid var(--border-medium); background: var(--bg-primary); border-radius: 4px; cursor: pointer;">Clear Filters</button>' : 
                        ''
                    }
                </div>
            `;
            elements.adrList.appendChild(noResults);
            return;
        }

        state.filteredAdrs.forEach(adr => {
            const listItem = document.createElement('li');
            listItem.className = 'adr-item';
            
            const link = document.createElement('a');
            link.href = `#adr-${adr.number}`;
            link.className = 'adr-link';
            link.setAttribute('data-adr', adr.number);
            
            link.innerHTML = `
                <div class="adr-header">
                    <span class="adr-number">ADR-${adr.number}</span>
                    <span class="adr-status ${adr.status.toLowerCase()}" data-status="${adr.status}"></span>
                </div>
                <div class="adr-title">${adr.title}</div>
            `;
            
            link.addEventListener('click', (e) => {
                e.preventDefault();
                loadADR(adr.number);
                closeMobileMenu();
            });
            
            listItem.appendChild(link);
            elements.adrList.appendChild(listItem);
        });
    }

    /**
     * Render ADR content in the main area
     */
    function renderADRContent(adr) {
        if (!elements.adrDocument || !adr.content) return;

        // Hide welcome screen, show content
        if (elements.welcomeScreen) {
            elements.welcomeScreen.style.display = 'none';
        }
        if (elements.errorScreen) {
            elements.errorScreen.style.display = 'none';
        }
        if (elements.adrContent) {
            elements.adrContent.style.display = 'block';
        }

        // Generate table of contents first
        const tocContent = generateTableOfContents(adr.content);
        
        // Render markdown content
        const htmlContent = ADRRenderer.renderMarkdown(adr.content, adr);
        
        // Create content with TOC layout
        if (tocContent && tocContent.length > 0) {
            elements.adrContent.innerHTML = `
                <div class="content-with-toc">
                    <nav class="table-of-contents">
                        <h3>Contents</h3>
                        <ul id="toc-list">${tocContent}</ul>
                    </nav>
                    <div class="adr-document">${htmlContent}</div>
                </div>
            `;
        } else {
            elements.adrContent.innerHTML = `<div class="adr-document no-toc">${htmlContent}</div>`;
        }

        // Update TOC references
        elements.adrDocument = elements.adrContent.querySelector('.adr-document');
        elements.tocList = elements.adrContent.querySelector('#toc-list');

        // Add heading IDs and TOC functionality
        setupTableOfContents();

        // Process mermaid diagrams
        processMermaidDiagrams();

        // Re-run Prism syntax highlighting for dynamically loaded content
        if (window.Prism) {
            setTimeout(() => {
                window.Prism.highlightAllUnder(elements.adrDocument);
            }, 0);
        }

        // Scroll to top
        elements.mainContent.scrollTop = 0;
    }

    /**
     * Show welcome screen
     */
    function showWelcomeScreen() {
        if (elements.welcomeScreen) {
            elements.welcomeScreen.style.display = 'block';
        }
        if (elements.adrContent) {
            elements.adrContent.style.display = 'none';
        }
        if (elements.errorScreen) {
            elements.errorScreen.style.display = 'none';
        }
        
        updateBreadcrumb(null);
        updateActiveNavigation(null);
        
        // Clear URL hash
        if (window.location.hash) {
            window.history.pushState(null, null, window.location.pathname);
        }
    }

    /**
     * Show error screen
     */
    function showError(message) {
        if (elements.errorScreen) {
            const errorMessage = elements.errorScreen.querySelector('#error-message');
            if (errorMessage) {
                errorMessage.textContent = message;
            }
            elements.errorScreen.style.display = 'flex';
        }
        
        if (elements.welcomeScreen) {
            elements.welcomeScreen.style.display = 'none';
        }
        if (elements.adrContent) {
            elements.adrContent.style.display = 'none';
        }

        // Bind retry button
        const retryBtn = document.getElementById('retry-btn');
        if (retryBtn) {
            retryBtn.onclick = () => {
                showWelcomeScreen();
                loadADRs();
            };
        }
    }

    /**
     * Generate table of contents from markdown content
     */
    function generateTableOfContents(content) {
        if (!content) return '';

        const lines = content.split('\n');
        const headings = [];
        
        lines.forEach((line, index) => {
            const headingMatch = line.match(/^(#{2,4})\s+(.+)$/);
            if (headingMatch) {
                const level = headingMatch[1].length;
                const text = headingMatch[2].trim();
                // Create slug-style ID from heading text
                const id = text.toLowerCase()
                    .replace(/[^\w\s-]/g, '') // Remove special characters
                    .replace(/\s+/g, '-') // Replace spaces with hyphens
                    .trim();
                
                headings.push({
                    level,
                    text,
                    id: id || `heading-${index}` // Fallback to numbered ID
                });
            }
        });

        if (headings.length === 0) return '';

        return headings.map(heading => 
            `<li><a href="#${heading.id}" class="toc-level-h${heading.level}" data-heading-id="${heading.id}">${heading.text}</a></li>`
        ).join('');
    }

    /**
     * Setup table of contents functionality
     */
    function setupTableOfContents() {
        if (!elements.adrDocument) return;

        const headings = elements.adrDocument.querySelectorAll('h2, h3, h4');
        const tocLinks = elements.adrContent.querySelectorAll('.table-of-contents a');
        
        // Create a mapping of TOC links to their intended heading text
        const linkToHeadingMap = new Map();
        tocLinks.forEach(link => {
            const headingText = link.textContent.trim();
            linkToHeadingMap.set(link, headingText);
        });

        // Add IDs to headings based on their text content
        headings.forEach((heading, index) => {
            const headingText = heading.textContent.trim();
            
            // Create slug-style ID from heading text (matching TOC generation)
            let id = headingText.toLowerCase()
                .replace(/[^\w\s-]/g, '') // Remove special characters
                .replace(/\s+/g, '-') // Replace spaces with hyphens
                .trim();
            
            // Fallback to numbered ID if slug is empty
            if (!id) {
                id = `heading-${index}`;
            }
            
            heading.id = id;
        });

        // Add click handlers to TOC links
        tocLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                
                const headingId = link.getAttribute('data-heading-id');
                let heading = document.getElementById(headingId);
                
                // If heading not found by ID, try to find by text content
                if (!heading) {
                    const headingText = link.textContent.trim();
                    heading = Array.from(headings).find(h => 
                        h.textContent.trim() === headingText
                    );
                }
                
                if (heading) {
                    // Remove active class from all TOC links
                    tocLinks.forEach(l => l.classList.remove('active'));
                    // Add active class to clicked link
                    link.classList.add('active');
                    
                    // Calculate scroll position to account for fixed header
                    const headerHeight = 60; // Approximate header height
                    const elementPosition = heading.getBoundingClientRect().top;
                    const offsetPosition = elementPosition + window.pageYOffset - headerHeight;
                    
                    // Smooth scroll to position
                    window.scrollTo({
                        top: offsetPosition,
                        behavior: 'smooth'
                    });
                }
            });
        });

        // Update active TOC link on scroll
        if (tocLinks.length > 0 && headings.length > 0) {
            setupTocScrollTracking(tocLinks, headings);
        }
    }

    /**
     * Setup scroll tracking for TOC
     */
    function setupTocScrollTracking(tocLinks, headings) {
        const observer = new IntersectionObserver((entries) => {
            let activeHeading = null;
            
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    activeHeading = entry.target;
                }
            });

            if (activeHeading) {
                // Remove active class from all links
                tocLinks.forEach(link => link.classList.remove('active'));
                
                // Find corresponding TOC link
                const activeLink = Array.from(tocLinks).find(link => 
                    link.getAttribute('data-heading-id') === activeHeading.id
                );
                
                if (activeLink) {
                    activeLink.classList.add('active');
                }
            }
        }, {
            rootMargin: '-20% 0px -60% 0px',
            threshold: 0
        });

        headings.forEach(heading => observer.observe(heading));
    }

    /**
     * Process Mermaid diagrams
     */
    async function processMermaidDiagrams() {
        const mermaidBlocks = elements.adrDocument.querySelectorAll('pre code.language-mermaid, .mermaid');
        
        for (let i = 0; i < mermaidBlocks.length; i++) {
            const block = mermaidBlocks[i];
            const diagramText = block.textContent;
            
            try {
                const diagramId = `mermaid-${Date.now()}-${i}`;
                const { svg } = await mermaid.render(diagramId, diagramText);
                
                // Create container
                const container = document.createElement('div');
                container.className = 'mermaid-container';
                container.innerHTML = `
                    ${svg}
                    <div class="diagram-actions">
                        <button class="diagram-btn" onclick="ADRApp.openDiagramModal(this)" title="View fullscreen">⛶</button>
                        <button class="diagram-btn" onclick="ADRApp.copyDiagram(this)" title="Copy SVG">📋</button>
                    </div>
                `;
                
                // Replace the original block
                block.closest('pre')?.replaceWith(container) || block.replaceWith(container);
                
            } catch (error) {
                console.error('Failed to render Mermaid diagram:', error);
                
                // Show error in place of diagram
                const errorDiv = document.createElement('div');
                errorDiv.className = 'mermaid-error';
                errorDiv.innerHTML = `
                    <div style="padding: var(--spacing-lg); background: var(--bg-secondary); border: 1px solid var(--border-medium); border-radius: 6px; text-align: center;">
                        <p style="color: var(--text-muted); margin: 0;">Failed to render diagram</p>
                        <details style="margin-top: var(--spacing-sm);">
                            <summary style="cursor: pointer; color: var(--text-secondary);">View source</summary>
                            <pre style="text-align: left; margin-top: var(--spacing-sm); font-size: 12px;"><code>${diagramText}</code></pre>
                        </details>
                    </div>
                `;
                
                block.closest('pre')?.replaceWith(errorDiv) || block.replaceWith(errorDiv);
            }
        }
    }

    /**
     * Update statistics
     */
    function updateStatistics() {
        stats.totalAdrs = state.adrs.length;
        stats.acceptedAdrs = state.adrs.filter(adr => adr.status.toLowerCase() === 'accepted').length;
        stats.c4Diagrams = state.adrs.filter(adr => adr.diagramType && adr.diagramType !== '-').length;

        // Update DOM
        const totalElement = document.getElementById('total-adrs');
        const acceptedElement = document.getElementById('accepted-adrs');
        const diagramsElement = document.getElementById('c4-diagrams');

        if (totalElement) totalElement.textContent = stats.totalAdrs;
        if (acceptedElement) acceptedElement.textContent = stats.acceptedAdrs;
        if (diagramsElement) diagramsElement.textContent = stats.c4Diagrams;
    }

    /**
     * Update active navigation
     */
    function updateActiveNavigation(adrNumber) {
        // Remove active class from all links
        const allLinks = elements.adrList.querySelectorAll('.adr-link');
        allLinks.forEach(link => link.classList.remove('active'));

        // Add active class to current link
        if (adrNumber) {
            const currentLink = elements.adrList.querySelector(`[data-adr="${adrNumber}"]`);
            if (currentLink) {
                currentLink.classList.add('active');
                currentLink.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
            }
        }
    }

    /**
     * Update breadcrumb
     */
    function updateBreadcrumb(adr) {
        if (!elements.currentAdrBreadcrumb) return;

        if (adr) {
            elements.currentAdrBreadcrumb.textContent = `ADR-${adr.number}: ${adr.title}`;
        } else {
            elements.currentAdrBreadcrumb.textContent = '';
        }
    }

    /**
     * Handle search input
     */
    function handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        state.searchTerm = searchTerm;

        // Show/hide clear button
        if (elements.searchClear) {
            elements.searchClear.classList.toggle('visible', searchTerm.length > 0);
        }

        filterADRs();
    }

    /**
     * Handle search keydown
     */
    function handleSearchKeydown(e) {
        if (e.key === 'Escape') {
            clearSearch();
        } else if (e.key === 'Enter') {
            // Navigate to first result
            if (state.filteredAdrs.length > 0) {
                loadADR(state.filteredAdrs[0].number);
                closeMobileMenu();
            }
        }
    }

    /**
     * Clear search
     */
    function clearSearch() {
        if (elements.searchInput) {
            elements.searchInput.value = '';
        }
        if (elements.searchClear) {
            elements.searchClear.classList.remove('visible');
        }
        state.searchTerm = '';
        filterADRs();
    }

    /**
     * Handle status filter
     */
    function handleStatusFilter(e) {
        state.statusFilter = e.target.value;
        filterADRs();
    }

    /**
     * Filter ADRs based on search and status
     */
    function filterADRs() {
        state.filteredAdrs = state.adrs.filter(adr => {
            const matchesSearch = !state.searchTerm || 
                adr.title.toLowerCase().includes(state.searchTerm) ||
                adr.number.includes(state.searchTerm) ||
                adr.status.toLowerCase().includes(state.searchTerm);

            const matchesStatus = !state.statusFilter || 
                adr.status === state.statusFilter;

            return matchesSearch && matchesStatus;
        });

        renderADRList();
    }

    /**
     * Clear all filters
     */
    function clearFilters() {
        clearSearch();
        if (elements.statusFilter) {
            elements.statusFilter.value = '';
        }
        state.statusFilter = '';
        filterADRs();
    }

    /**
     * Toggle mobile menu
     */
    function toggleMobileMenu() {
        state.isMobileMenuOpen = !state.isMobileMenuOpen;
        
        if (elements.sidebar) {
            elements.sidebar.classList.toggle('open', state.isMobileMenuOpen);
        }
        if (elements.mobileOverlay) {
            elements.mobileOverlay.classList.toggle('active', state.isMobileMenuOpen);
        }
    }

    /**
     * Close mobile menu
     */
    function closeMobileMenu() {
        state.isMobileMenuOpen = false;
        
        if (elements.sidebar) {
            elements.sidebar.classList.remove('open');
        }
        if (elements.mobileOverlay) {
            elements.mobileOverlay.classList.remove('active');
        }
    }

    /**
     * Handle hash change
     */
    function handleHashChange() {
        const hash = window.location.hash;
        
        if (hash && hash.startsWith('#adr-')) {
            const adrNumber = hash.replace('#adr-', '');
            if (adrNumber !== state.currentAdr?.number) {
                loadADR(adrNumber);
            }
        } else if (hash === '#docs') {
            if (state.currentAdr?.number !== 'DOCS') {
                showDocumentation();
            }
        } else if (hash === '#template') {
            if (state.currentAdr?.number !== 'TEMPLATE') {
                showTemplate();
            }
        } else if (hash === '' && state.currentAdr) {
            showWelcomeScreen();
        }
    }

    /**
     * Handle window resize
     */
    function handleResize() {
        // Close mobile menu on desktop
        if (window.innerWidth > 768 && state.isMobileMenuOpen) {
            closeMobileMenu();
        }
    }

    /**
     * Handle keyboard shortcuts
     */
    function handleKeyboardShortcuts(e) {
        // Cmd/Ctrl + K to focus search
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
            e.preventDefault();
            if (elements.searchInput) {
                elements.searchInput.focus();
            }
        }

        // Escape to close mobile menu or clear search
        if (e.key === 'Escape') {
            if (state.isMobileMenuOpen) {
                closeMobileMenu();
            } else if (state.searchTerm) {
                clearSearch();
            }
        }
    }

    /**
     * Print current ADR
     */
    function printCurrentADR() {
        if (state.currentAdr) {
            window.print();
        }
    }

    /**
     * Toggle fullscreen
     */
    function toggleFullscreen() {
        if (!document.fullscreenElement) {
            document.documentElement.requestFullscreen();
        } else {
            document.exitFullscreen();
        }
    }

    /**
     * Navigate to home page
     */
    function goHome() {
        // Clear any current ADR
        state.currentAdr = null;
        
        // Show welcome screen
        showWelcomeScreen();
        
        // Close mobile menu if open
        closeMobileMenu();
        
        // Clear URL hash
        if (window.location.hash) {
            window.history.pushState(null, null, window.location.pathname);
        }
    }

    /**
     * Show documentation
     */
    async function showDocumentation() {
        setLoading(true);
        
        try {
            // Create a documentation ADR object
            const docContent = generateDocumentationContent();
            const docAdr = {
                number: 'DOCS',
                title: 'ADR Documentation & Guide',
                status: 'Reference',
                content: docContent
            };
            
            state.currentAdr = docAdr;
            renderADRContent(docAdr);
            updateActiveNavigation(null); // Don't highlight any sidebar item
            updateBreadcrumb({ number: 'DOCS', title: 'Documentation' });
            closeMobileMenu();
            
            // Update URL hash
            window.history.pushState(null, null, '#docs');
            
        } catch (error) {
            console.error('Failed to show documentation:', error);
            showError('Failed to load documentation');
        } finally {
            setLoading(false);
        }
    }

    /**
     * Show ADR template
     */
    async function showTemplate() {
        setLoading(true);
        
        try {
            const response = await fetch('adr/template.md');
            if (!response.ok) {
                throw new Error('Template not found');
            }
            
            const content = await response.text();
            const templateAdr = {
                number: 'TEMPLATE',
                title: 'ADR Template',
                status: 'Template',
                content: content
            };
            
            state.currentAdr = templateAdr;
            renderADRContent(templateAdr);
            updateActiveNavigation(null);
            updateBreadcrumb({ number: 'TEMPLATE', title: 'ADR Template' });
            closeMobileMenu();
            
            // Update URL hash
            window.history.pushState(null, null, '#template');
            
        } catch (error) {
            console.error('Failed to load template:', error);
            showError('Failed to load ADR template');
        } finally {
            setLoading(false);
        }
    }

    /**
     * Generate documentation content
     */
    function generateDocumentationContent() {
        return `# ADR Documentation & Guide

## What are Architecture Decision Records (ADRs)?

Architecture Decision Records (ADRs) are **lightweight documents that capture important architectural decisions** made during a project's development. Each ADR documents:

- **The decision made** and its rationale
- **The context** that led to the decision  
- **The consequences** (both positive and negative) of the decision
- **Alternative options** that were considered

ADRs help teams understand the reasoning behind architectural choices and maintain a historical record of decision-making that survives team changes and time.

## Why Use ADRs?

### Key Benefits

- **Knowledge Preservation**: Architectural decisions and their reasoning are documented for future reference
- **Onboarding**: New team members can understand the system's evolution and current state
- **Decision Transparency**: Everyone can see what decisions were made and why
- **Avoid Repeated Discussions**: Settled architectural matters don't need to be re-debated
- **Change Management**: Understanding past decisions helps evaluate future changes
- **Accountability**: Clear record of who decided what and when

### When to Create an ADR

Create an ADR when making decisions that:
- Affect the overall system architecture
- Have long-term consequences
- Involve significant trade-offs
- Impact multiple teams or components
- Introduce new technologies or patterns
- Change existing architectural patterns

## ADR Status Lifecycle

### Status Types

- **⏳ Proposed** - Under review and discussion
- **✅ Accepted** - Approved and being implemented
- **🗑️ Deprecated** - No longer recommended but still in use
- **⬆️ Superseded** - Replaced by a newer decision

### Decision Process

1. **Identify Need** - Recognize when an architectural decision is required
2. **Research & Analysis** - Gather information and evaluate options
3. **Draft ADR** - Create the initial document with all sections
4. **Review Process** - Team review and Architecture Review Board approval
5. **Decision** - Final status determination (Accepted/Rejected)
6. **Implementation** - Execute the decision
7. **Maintenance** - Update status as needed over time

## ADR Structure

Each ADR follows a standard format:

### Required Sections

- **Title** - Clear, concise description of the decision
- **Status** - Current state of the decision
- **Context** - Background and forces that led to this decision
- **Decision** - The choice that was made and why
- **Consequences** - Both positive and negative outcomes

### Optional Sections

- **Alternatives Considered** - Other options that were evaluated
- **Implementation Notes** - Specific guidance for implementation
- **Related Decisions** - Links to other relevant ADRs

## Best Practices

### Writing Effective ADRs

1. **Be Concise but Complete** - Capture essential information without unnecessary detail
2. **Use Plain Language** - Avoid jargon; write for future team members
3. **Include Visuals** - Use diagrams to clarify complex decisions
4. **Show Your Work** - Explain alternatives and trade-offs considered
5. **Be Honest** - Include negative consequences and risks
6. **Keep It Current** - Update status when decisions change

### C4 Model Integration

This repository enhances ADRs with C4 model diagrams using Mermaid syntax:

- **Context Diagrams** - Show system boundaries and external dependencies
- **Container Diagrams** - Show high-level technology choices and data flow
- **Component Diagrams** - Show internal structure of containers/services
- **Dynamic Diagrams** - Show workflow and process flows over time

## Using This Interface

### Navigation
- **Sidebar** - Browse all ADRs, search, and filter by status
- **Home Button** - Click the title or icon to return to the welcome screen
- **Table of Contents** - Navigate within long ADR documents
- **Breadcrumbs** - See your current location and navigate back

### Features
- **Search** - Find ADRs by title, number, or content
- **Filtering** - View ADRs by status (Proposed, Accepted, etc.)
- **Dark Mode** - Automatically follows your system preference
- **Responsive Design** - Works on desktop, tablet, and mobile
- **Printable** - Use the print button to create PDFs

### Interactive Diagrams
- **Full Screen** - Click the expand button on diagrams
- **Copy SVG** - Use the copy button to export diagram code
- **Mermaid Integration** - All diagrams are rendered from markdown

## Contributing

To add a new ADR:

1. **Use the Template** - Start with the ADR template
2. **Sequential Numbering** - Use the next available number (0001, 0002, etc.)
3. **Follow Format** - Include all required sections
4. **Add Diagrams** - Include C4 diagrams where helpful
5. **Review Process** - Submit for team and ARB review
6. **Update Index** - The system automatically updates the index

### Naming Convention

ADR files should follow this pattern:
\`NNNN-kebab-case-title.md\`

Examples:
- \`0001-record-architecture-decisions.md\`
- \`0002-establish-architecture-review-board.md\`
- \`0003-adopt-microservices-architecture.md\`

## Architecture Review Board (ARB)

Our Architecture Review Board ensures decisions are:
- **Technically Sound** - Feasible and well-reasoned
- **Aligned** - Consistent with existing architecture
- **Business Valuable** - Support business objectives
- **Risk Aware** - Consider and mitigate potential issues

## Tools & Automation

This repository includes:
- **GitHub Actions** - Automated validation and deployment
- **Makefile** - Local development commands
- **Web Interface** - This documentation and ADR browser
- **Index Generation** - Automatic ADR discovery and indexing

### Development Commands

\`\`\`bash
make generate-index  # Generate ADR index
make serve          # Start local development server
make validate       # Validate ADR structure
make new-adr        # Create new ADR from template
\`\`\`

---

This documentation is part of the ShopFlow ADR repository, demonstrating best practices for documenting and managing architectural decisions.`;
    }

    /**
     * Set loading state
     */
    function setLoading(loading) {
        state.isLoading = loading;
        
        if (elements.loading) {
            elements.loading.classList.toggle('hidden', !loading);
        }
    }

    /**
     * Open diagram modal
     */
    function openDiagramModal(button) {
        const container = button.closest('.mermaid-container');
        const svg = container.querySelector('svg');
        
        if (!svg) return;

        // Create modal
        const modal = document.createElement('div');
        modal.className = 'diagram-modal active';
        modal.innerHTML = `
            <div class="diagram-modal-content">
                <button class="diagram-modal-close" onclick="this.closest('.diagram-modal').remove()">×</button>
                ${svg.outerHTML}
            </div>
        `;

        document.body.appendChild(modal);

        // Close on overlay click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });

        // Close on escape
        const handleEscape = (e) => {
            if (e.key === 'Escape') {
                modal.remove();
                document.removeEventListener('keydown', handleEscape);
            }
        };
        document.addEventListener('keydown', handleEscape);
    }

    /**
     * Copy diagram SVG
     */
    async function copyDiagram(button) {
        const container = button.closest('.mermaid-container');
        const svg = container.querySelector('svg');
        
        if (!svg) return;

        try {
            await navigator.clipboard.writeText(svg.outerHTML);
            
            // Show feedback
            const originalText = button.textContent;
            button.textContent = '✓';
            button.style.background = 'var(--status-accepted)';
            button.style.color = 'white';
            
            setTimeout(() => {
                button.textContent = originalText;
                button.style.background = '';
                button.style.color = '';
            }, 1500);
            
        } catch (error) {
            console.error('Failed to copy diagram:', error);
            
            // Show error feedback
            const originalText = button.textContent;
            button.textContent = '✗';
            button.style.background = 'var(--status-deprecated)';
            button.style.color = 'white';
            
            setTimeout(() => {
                button.textContent = originalText;
                button.style.background = '';
                button.style.color = '';
            }, 1500);
        }
    }

    // Public API
    return {
        init,
        loadADR,
        clearFilters,
        openDiagramModal,
        copyDiagram,
        showDocumentation,
        showTemplate,
        goHome,
        
        // State getters
        getCurrentADR: () => state.currentAdr,
        getADRs: () => state.adrs,
        getStats: () => stats
    };

})();