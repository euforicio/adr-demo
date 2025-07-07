// ADR Generator - Enhanced JavaScript
console.log('ADR site loaded with enhanced features');


// Mobile sidebar toggle
function initializeMobileSidebar() {
    const sidebar = document.querySelector('.sidebar');
    const mainContent = document.querySelector('.main-content');
    
    if (!sidebar || window.innerWidth > 768) return;
    
    // Create mobile menu button
    const menuButton = document.createElement('button');
    menuButton.className = 'mobile-menu-btn';
    menuButton.innerHTML = '☰';
    menuButton.style.cssText = `
        position: fixed;
        top: 1rem;
        left: 1rem;
        z-index: 1001;
        background: #5e6ad2;
        color: white;
        border: none;
        padding: 0.5rem;
        border-radius: 0.375rem;
        font-size: 1.25rem;
        cursor: pointer;
        display: block;
    `;
    
    document.body.appendChild(menuButton);
    
    menuButton.addEventListener('click', function() {
        sidebar.classList.toggle('open');
    });
    
    // Close sidebar when clicking outside
    document.addEventListener('click', function(e) {
        if (!sidebar.contains(e.target) && !menuButton.contains(e.target)) {
            sidebar.classList.remove('open');
        }
    });
}

// Enhanced keyboard navigation
function initializeKeyboardNavigation() {
    document.addEventListener('keydown', function(e) {
        // Ctrl/Cmd + K to focus search
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
            e.preventDefault();
            const searchInput = document.getElementById('search-input');
            if (searchInput) {
                searchInput.focus();
                searchInput.select();
            }
        }
        
        // Escape to clear search
        if (e.key === 'Escape') {
            const searchInput = document.getElementById('search-input');
            const searchClear = document.getElementById('search-clear');
            if (searchInput && searchInput.value) {
                if (searchClear) {
                    searchClear.click();
                } else {
                    searchInput.value = '';
                    searchInput.dispatchEvent(new Event('input'));
                }
            }
        }
    });
}

// Smooth scrolling for anchor links
function initializeSmoothScrolling() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
}

// Initialize mermaid diagrams
function initializeMermaid() {
    // Mermaid should be initialized by the mermaid script in the template
    // This function handles post-processing and interactions
    
    const mermaidElements = document.querySelectorAll('.mermaid');
    
    mermaidElements.forEach((element, index) => {
        // Add click to enlarge functionality
        element.style.cursor = 'pointer';
        element.title = 'Click to view fullscreen';
        
        element.addEventListener('click', function() {
            const svg = element.querySelector('svg');
            if (svg) {
                openMermaidFullscreen(element, svg);
            }
        });
        
        // Store original diagram code for copying
        const preElement = element.previousElementSibling;
        if (preElement && preElement.tagName === 'PRE') {
            const code = preElement.textContent;
            element.setAttribute('data-diagram', code);
        }
    });
}

// Open mermaid diagram in fullscreen modal
function openMermaidFullscreen(diagramElement, svg) {
    const modal = document.createElement('div');
    modal.className = 'mermaid-modal';
    
    const modalContent = document.createElement('div');
    modalContent.className = 'mermaid-modal-content';
    
    const closeButton = document.createElement('button');
    closeButton.className = 'mermaid-modal-close';
    closeButton.innerHTML = '×';
    closeButton.title = 'Close';
    closeButton.addEventListener('click', () => modal.remove());
    
    const clonedSvg = svg.cloneNode(true);
    clonedSvg.style.maxWidth = '100%';
    clonedSvg.style.height = 'auto';
    
    modalContent.appendChild(closeButton);
    modalContent.appendChild(clonedSvg);
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    // Close on escape key
    document.addEventListener('keydown', function escapeHandler(e) {
        if (e.key === 'Escape') {
            modal.remove();
            document.removeEventListener('keydown', escapeHandler);
        }
    });
    
    // Close on backdrop click
    modal.addEventListener('click', function(e) {
        if (e.target === modal) {
            modal.remove();
        }
    });
}

// Copy mermaid diagram code to clipboard
function copyMermaidCode(diagramElement) {
    const code = diagramElement.getAttribute('data-diagram');
    if (code && navigator.clipboard) {
        navigator.clipboard.writeText(code).then(() => {
            // Show temporary success message
            const button = event.target;
            const originalText = button.textContent;
            button.textContent = 'Copied!';
            button.style.background = '#10b981';
            setTimeout(() => {
                button.textContent = originalText;
                button.style.background = '';
            }, 2000);
        }).catch(err => {
            console.error('Failed to copy diagram code:', err);
        });
    }
}

// Process code blocks for syntax highlighting hints
function processCodeBlocks() {
    const codeBlocks = document.querySelectorAll('pre code[class*="language-"]');
    
    codeBlocks.forEach(block => {
        const pre = block.parentElement;
        const language = block.className.match(/language-(\w+)/);
        
        if (language) {
            pre.setAttribute('data-language', language[1]);
            pre.className = `language-${language[1]}`;
        }
    });
}

// ADR List View Toggle
function initializeViewToggle() {
    console.log('Initializing view toggle...');
    
    const groupedBtn = document.getElementById('view-grouped');
    const flatBtn = document.getElementById('view-flat');
    const groupedView = document.getElementById('adr-list-grouped');
    const flatView = document.getElementById('adr-list-flat');
    
    console.log('Toggle elements:', { groupedBtn, flatBtn, groupedView, flatView });
    
    if (!groupedBtn || !flatBtn || !groupedView || !flatView) {
        console.error('Missing toggle elements');
        return;
    }
    
    // Load saved preference or default to grouped
    const savedView = localStorage.getItem('adr-view-preference') || 'grouped';
    console.log('Saved view preference:', savedView);
    
    function setActiveView(viewType) {
        console.log('Setting active view to:', viewType);
        
        if (viewType === 'grouped') {
            // Show grouped view
            groupedView.classList.remove('hidden');
            flatView.classList.add('hidden');
            
            // Update button styles
            groupedBtn.classList.add('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-gray-100', 'shadow-sm');
            groupedBtn.classList.remove('text-gray-600', 'dark:text-gray-400');
            
            flatBtn.classList.remove('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-gray-100', 'shadow-sm');
            flatBtn.classList.add('text-gray-600', 'dark:text-gray-400');
        } else {
            // Show flat view
            groupedView.classList.add('hidden');
            flatView.classList.remove('hidden');
            
            // Update button styles
            flatBtn.classList.add('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-gray-100', 'shadow-sm');
            flatBtn.classList.remove('text-gray-600', 'dark:text-gray-400');
            
            groupedBtn.classList.remove('bg-white', 'dark:bg-gray-600', 'text-gray-900', 'dark:text-gray-100', 'shadow-sm');
            groupedBtn.classList.add('text-gray-600', 'dark:text-gray-400');
        }
        
        // Save preference
        localStorage.setItem('adr-view-preference', viewType);
        console.log('View preference saved:', viewType);
        
        // Re-run search to update visibility for new view
        setTimeout(() => {
            const searchInput = document.getElementById('search-input');
            if (searchInput) {
                searchInput.dispatchEvent(new Event('input'));
            }
        }, 100);
    }
    
    // Set initial view
    setActiveView(savedView);
    
    // Add click listeners with debugging
    groupedBtn.addEventListener('click', (e) => {
        console.log('Grouped button clicked');
        e.preventDefault();
        e.stopPropagation();
        setActiveView('grouped');
    });
    
    flatBtn.addEventListener('click', (e) => {
        console.log('Flat button clicked');
        e.preventDefault();
        e.stopPropagation();
        setActiveView('flat');
    });
    
    console.log('Event listeners attached to buttons');
    
    // Test button accessibility
    console.log('Button states:', {
        groupedDisabled: groupedBtn.disabled,
        flatDisabled: flatBtn.disabled,
        groupedVisible: window.getComputedStyle(groupedBtn).display !== 'none',
        flatVisible: window.getComputedStyle(flatBtn).display !== 'none'
    });
    
    console.log('View toggle initialized successfully');
}

// Enhanced search functionality to work with both views
function initializeEnhancedSearch() {
    const searchInput = document.getElementById('search-input');
    const searchClear = document.getElementById('search-clear');
    const statusFilter = document.getElementById('status-filter');
    
    if (!searchInput) return;
    
    function performSearch() {
        const query = searchInput.value.toLowerCase();
        const selectedStatus = statusFilter ? statusFilter.value.toLowerCase() : '';
        
        // Search in both grouped and flat views
        const groupedItems = document.querySelectorAll('#adr-list-grouped li');
        const flatItems = document.querySelectorAll('#adr-list-flat li');
        
        function filterItems(items) {
            items.forEach(item => {
                const text = item.textContent.toLowerCase();
                const statusIcon = item.querySelector('.w-5.h-5 span, .inline-flex'); // Status icon or badge
                const statusText = getStatusFromElement(statusIcon);
                
                const matchesQuery = !query || text.includes(query);
                const matchesStatus = !selectedStatus || statusText === selectedStatus;
                
                item.style.display = (matchesQuery && matchesStatus) ? '' : 'none';
            });
        }
        
        filterItems(groupedItems);
        filterItems(flatItems);
        
        // Hide empty categories in grouped view
        const categories = document.querySelectorAll('#adr-list-grouped > div');
        categories.forEach(category => {
            const visibleItems = Array.from(category.querySelectorAll('li')).filter(
                item => item.style.display !== 'none'
            );
            category.style.display = visibleItems.length > 0 ? '' : 'none';
        });
    }
    
    function getStatusFromElement(element) {
        if (!element) return '';
        const text = element.textContent.toLowerCase();
        
        if (text.includes('accepted') || text.includes('✓')) return 'accepted';
        if (text.includes('proposed') || text.includes('●')) return 'proposed';
        if (text.includes('deprecated') || text.includes('✗')) return 'deprecated';
        if (text.includes('superseded') || text.includes('↑')) return 'superseded';
        
        return '';
    }
    
    searchInput.addEventListener('input', performSearch);
    
    if (statusFilter) {
        statusFilter.addEventListener('change', performSearch);
    }
    
    if (searchClear) {
        searchClear.addEventListener('click', function() {
            searchInput.value = '';
            if (statusFilter) statusFilter.value = '';
            performSearch();
            searchInput.focus();
        });
    }
}

// Initialize all features when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing features...');
    initializeViewToggle();
    initializeEnhancedSearch();
    initializeMobileSidebar();
    initializeKeyboardNavigation();
    initializeSmoothScrolling();
    processCodeBlocks();
    initializeCategories();
    
    // Initialize mermaid after a short delay to ensure it's loaded
    setTimeout(initializeMermaid, 100);
    
    console.log('All ADR site features initialized');
});

// Handle window resize for mobile sidebar
window.addEventListener('resize', function() {
    if (window.innerWidth > 768) {
        const sidebar = document.querySelector('.sidebar');
        const menuButton = document.querySelector('.mobile-menu-btn');
        if (sidebar) sidebar.classList.remove('open');
        if (menuButton) menuButton.style.display = 'none';
    } else {
        const menuButton = document.querySelector('.mobile-menu-btn');
        if (menuButton) menuButton.style.display = 'block';
    }
});

// Category toggle function for hierarchical navigation
function toggleCategory(categoryId) {
    console.log('Toggle category called with ID:', categoryId);
    
    const category = document.getElementById(categoryId);
    const button = document.querySelector(`[onclick="toggleCategory('${categoryId}')"]`);
    const icon = button ? button.querySelector('.category-icon') : null;
    
    console.log('Found elements:', { category, button, icon });
    
    if (category && button && icon) {
        const isCollapsed = category.classList.contains('collapsed');
        console.log('Current collapsed state:', isCollapsed);
        
        if (isCollapsed) {
            // Show the category
            category.classList.remove('collapsed', 'hidden');
            category.style.display = '';
            button.classList.remove('collapsed');
            icon.textContent = '▼';
            console.log('Expanded category', category.className);
        } else {
            // Hide the category
            category.classList.add('collapsed', 'hidden');
            category.style.display = 'none';
            button.classList.add('collapsed');
            icon.textContent = '▶';
            console.log('Collapsed category', category.className);
        }
        
        // Force a style recalculation
        category.offsetHeight;
    } else {
        console.error('Missing elements for category toggle:', { categoryId, category, button, icon });
    }
}

// Initialize category states on load
function initializeCategories() {
    console.log('Initializing categories...');
    
    // Set default state (all categories expanded)
    const categories = document.querySelectorAll('.adr-category-list');
    const buttons = document.querySelectorAll('.adr-category-header');
    
    console.log('Found categories:', categories.length, 'buttons:', buttons.length);
    
    categories.forEach(category => {
        category.classList.remove('collapsed');
    });
    
    buttons.forEach(button => {
        button.classList.remove('collapsed');
        const icon = button.querySelector('.category-icon');
        if (icon) {
            icon.textContent = '▼';
        }
    });
    
    console.log('Categories initialized as expanded');
}

// Global functions for template use
window.openMermaidFullscreen = openMermaidFullscreen;
window.copyMermaidCode = copyMermaidCode;
window.toggleCategory = toggleCategory;