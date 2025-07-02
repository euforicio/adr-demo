// ADR Renderer - Handles markdown rendering and ADR-specific formatting
window.ADRRenderer = (function() {
    'use strict';

    /**
     * Configure marked.js with custom settings
     */
    function configureMarked() {
        // Configure marked options for full GitHub Flavored Markdown support
        marked.setOptions({
            gfm: true,           // GitHub Flavored Markdown
            breaks: true,        // Line breaks
            pedantic: false,
            sanitize: false,
            smartLists: true,
            smartypants: true,
            xhtml: false,
            headerIds: true,     // Generate header IDs
            mangle: false,       // Don't mangle email addresses
            highlight: function(code, lang) {
                // Let Prism.js handle syntax highlighting
                if (lang && window.Prism && window.Prism.languages[lang]) {
                    try {
                        return window.Prism.highlight(code, window.Prism.languages[lang], lang);
                    } catch (e) {
                        console.warn('Prism highlighting failed:', e);
                    }
                }
                return code;
            }
        });

        // Custom renderer for ADR-specific elements
        const renderer = new marked.Renderer();

        // Custom heading renderer with proper IDs
        renderer.heading = function(text, level) {
            // Create slug-style ID from heading text
            const escapedText = text.toLowerCase()
                .replace(/[^\w\s-]/g, '') // Remove special characters
                .replace(/\s+/g, '-') // Replace spaces with hyphens
                .trim();
            const sectionClass = getSectionClass(text);
            
            return `<h${level} id="${escapedText}" class="adr-heading ${sectionClass}">
                ${text}
            </h${level}>`;
        };

        // Custom code block renderer for Mermaid diagrams
        renderer.code = function(code, language) {
            if (language === 'mermaid') {
                return `<pre class="mermaid-source"><code class="language-mermaid">${code}</code></pre>`;
            }

            // For other languages, use Prism.js
            if (language && window.Prism && window.Prism.languages[language]) {
                try {
                    const highlighted = window.Prism.highlight(code, window.Prism.languages[language], language);
                    return `<pre class="language-${language}"><code class="language-${language}">${highlighted}</code></pre>`;
                } catch (e) {
                    console.warn('Prism highlighting failed:', e);
                }
            }

            return `<pre><code class="language-${language || 'text'}">${code}</code></pre>`;
        };

        // Custom table renderer with enhanced styling
        renderer.table = function(header, body) {
            return `<div class="table-container">
                <table class="adr-table">
                    <thead>${header}</thead>
                    <tbody>${body}</tbody>
                </table>
            </div>`;
        };

        // Custom list item renderer for task lists (GitHub style checkboxes)
        renderer.listitem = function(text) {
            // Check for task list items
            if (/^\s*\[[x\s]\]\s*/.test(text)) {
                const isChecked = /^\s*\[x\]\s*/i.test(text);
                const cleanText = text.replace(/^\s*\[[x\s]\]\s*/i, '');
                const checkboxId = 'task-' + Math.random().toString(36).substr(2, 9);
                
                return `<li class="task-list-item">
                    <input type="checkbox" id="${checkboxId}" ${isChecked ? 'checked' : ''} disabled class="task-checkbox">
                    <label for="${checkboxId}" class="task-label">${cleanText}</label>
                </li>`;
            }
            return `<li>${text}</li>`;
        };

        // Enhanced image renderer with lazy loading and captions
        renderer.image = function(href, title, text) {
            const titleAttr = title ? ` title="${title}"` : '';
            const altText = text || '';
            
            return `<figure class="image-figure">
                <img src="${href}" alt="${altText}"${titleAttr} loading="lazy" class="adr-image">
                ${title ? `<figcaption class="image-caption">${title}</figcaption>` : ''}
            </figure>`;
        };

        // Enhanced horizontal rule renderer
        renderer.hr = function() {
            return '<hr class="adr-divider">';
        };

        // Custom list renderer
        renderer.list = function(body, ordered) {
            const type = ordered ? 'ol' : 'ul';
            const className = ordered ? 'adr-ordered-list' : 'adr-unordered-list';
            return `<${type} class="${className}">${body}</${type}>`;
        };

        // Custom blockquote renderer
        renderer.blockquote = function(quote) {
            return `<blockquote class="adr-blockquote">${quote}</blockquote>`;
        };

        // Custom link renderer with external link handling
        renderer.link = function(href, title, text) {
            const isExternal = href.startsWith('http') || href.startsWith('//');
            const titleAttr = title ? ` title="${title}"` : '';
            const target = isExternal ? ' target="_blank" rel="noopener noreferrer"' : '';
            const className = isExternal ? ' class="external-link"' : ' class="internal-link"';
            
            return `<a href="${href}"${titleAttr}${target}${className}>${text}</a>`;
        };

        marked.use({ renderer });
    }

    /**
     * Render markdown content to HTML
     */
    function renderMarkdown(content, adr) {
        if (!content) return '';

        try {
            // Configure marked if not already done
            if (!marked.configured) {
                configureMarked();
                marked.configured = true;
            }

            // Parse markdown to HTML
            let html = marked.parse(content);

            // Post-process HTML for ADR-specific enhancements
            html = enhanceADRSections(html);
            html = addADRNumber(html, adr);
            html = addStatusBadge(html, adr);
            html = enhanceConsequencesSection(html);
            html = addImplementationNotes(html);
            html = addRelatedADRs(html, adr);

            return html;

        } catch (error) {
            console.error('Failed to render markdown:', error);
            return `<div class="render-error">
                <h3>Rendering Error</h3>
                <p>Failed to render markdown content.</p>
                <details>
                    <summary>Error Details</summary>
                    <pre>${error.message}</pre>
                </details>
                <details>
                    <summary>Original Content</summary>
                    <pre>${content}</pre>
                </details>
            </div>`;
        }
    }

    /**
     * Get section class based on heading text
     */
    function getSectionClass(text) {
        const lowerText = text.toLowerCase();
        
        if (lowerText.includes('status')) return 'status';
        if (lowerText.includes('context')) return 'context';
        if (lowerText.includes('decision')) return 'decision';
        if (lowerText.includes('consequences')) return 'consequences';
        if (lowerText.includes('implementation')) return 'implementation';
        if (lowerText.includes('related')) return 'related';
        
        return '';
    }

    /**
     * Enhance ADR sections with better structure
     */
    function enhanceADRSections(html) {
        // Wrap sections in containers
        const sectionHeaders = ['Status', 'Context', 'Decision', 'Consequences'];
        
        sectionHeaders.forEach(section => {
            const regex = new RegExp(`(<h2[^>]*class="[^"]*${section.toLowerCase()}[^"]*"[^>]*>${section}</h2>)`, 'gi');
            html = html.replace(regex, (match) => {
                return `<div class="adr-section ${section.toLowerCase()}">${match}`;
            });
        });

        // Close section divs (simple approach - close before next h2 or at end)
        html = html.replace(/(<div class="adr-section[^"]*">[\s\S]*?)(?=<div class="adr-section|$)/g, '$1</div>');
        
        return html;
    }

    /**
     * Add ADR number to the main title
     */
    function addADRNumber(html, adr) {
        if (!adr || !adr.number) return html;

        // Find the main h1 title and prepend the ADR number
        const titleRegex = /(<h1[^>]*>)(.*?)(<\/h1>)/i;
        
        return html.replace(titleRegex, (match, opening, title, closing) => {
            // Don't add number if it's already there or for special docs
            if (title.includes('ADR-') || adr.number === 'DOCS' || adr.number === 'TEMPLATE') {
                return match;
            }
            
            return `${opening}<span class="adr-number">ADR-${adr.number}</span> ${title}${closing}`;
        });
    }

    /**
     * Add status badge to the status section
     */
    function addStatusBadge(html, adr) {
        if (!adr || !adr.status) return html;

        const statusRegex = /(<div class="adr-section status">[\s\S]*?<h2[^>]*>Status<\/h2>)([\s\S]*?)(<\/div>)/i;
        
        return html.replace(statusRegex, (match, opening, content, closing) => {
            const status = adr.status.toLowerCase();
            const badge = `<div class="status-badge ${status}">${adr.status}</div>`;
            
            return `${opening}${badge}${content}${closing}`;
        });
    }

    /**
     * Enhance consequences section with categorization
     */
    function enhanceConsequencesSection(html) {
        const consequencesRegex = /(<div class="adr-section consequences">[\s\S]*?)(<\/div>)/i;
        
        return html.replace(consequencesRegex, (match, section, closing) => {
            // Look for Positive/Negative/Neutral subsections
            let enhanced = section;
            
            // Enhance lists that follow common patterns
            enhanced = enhanced.replace(/(<h3[^>]*>Positive[^<]*<\/h3>[\s\S]*?)(<ul[^>]*>[\s\S]*?<\/ul>)/gi, 
                '$1<div class="consequence-category positive">$2</div>');
            
            enhanced = enhanced.replace(/(<h3[^>]*>Negative[^<]*<\/h3>[\s\S]*?)(<ul[^>]*>[\s\S]*?<\/ul>)/gi, 
                '$1<div class="consequence-category negative">$2</div>');
            
            enhanced = enhanced.replace(/(<h3[^>]*>Neutral[^<]*<\/h3>[\s\S]*?)(<ul[^>]*>[\s\S]*?<\/ul>)/gi, 
                '$1<div class="consequence-category neutral">$2</div>');

            // If no explicit categories, look for pattern indicators
            if (!enhanced.includes('consequence-category')) {
                enhanced = enhanceConsequencesByContent(enhanced);
            }

            return `${enhanced}${closing}`;
        });
    }

    /**
     * Enhance consequences by analyzing content patterns
     */
    function enhanceConsequencesByContent(sectionHtml) {
        // Split by lists and analyze content
        const listRegex = /<ul[^>]*>([\s\S]*?)<\/ul>/gi;
        
        return sectionHtml.replace(listRegex, (match, listContent) => {
            const items = listContent.match(/<li[^>]*>[\s\S]*?<\/li>/gi) || [];
            
            // Analyze sentiment of list items
            const positiveKeywords = ['improve', 'better', 'faster', 'easier', 'enable', 'support', 'enhance', 'optimize', 'scalable', 'flexible'];
            const negativeKeywords = ['complex', 'difficult', 'overhead', 'cost', 'risk', 'challenge', 'limitation', 'issue', 'problem', 'expensive'];
            
            let positiveScore = 0;
            let negativeScore = 0;
            
            items.forEach(item => {
                const text = item.toLowerCase();
                positiveKeywords.forEach(keyword => {
                    if (text.includes(keyword)) positiveScore++;
                });
                negativeKeywords.forEach(keyword => {
                    if (text.includes(keyword)) negativeScore++;
                });
            });

            let category = 'neutral';
            if (positiveScore > negativeScore) category = 'positive';
            else if (negativeScore > positiveScore) category = 'negative';

            return `<div class="consequence-category ${category}">${match}</div>`;
        });
    }

    /**
     * Add implementation notes section enhancement
     */
    function addImplementationNotes(html) {
        // Look for implementation-related sections
        const implementationRegex = /(<h3[^>]*>.*?[Ii]mplementation.*?<\/h3>)([\s\S]*?)(?=<h[1-6]|$)/gi;
        
        return html.replace(implementationRegex, (match, heading, content) => {
            return `<div class="implementation-notes">
                ${heading}
                ${content}
            </div>`;
        });
    }

    /**
     * Add related ADRs section
     */
    function addRelatedADRs(html, adr) {
        // This would be enhanced to actually find related ADRs
        // For now, just enhance any existing related sections
        const relatedRegex = /(<h3[^>]*>.*?[Rr]elated.*?<\/h3>)([\s\S]*?)(?=<h[1-6]|$)/gi;
        
        return html.replace(relatedRegex, (match, heading, content) => {
            return `<div class="related-adrs">
                ${heading}
                ${content}
            </div>`;
        });
    }

    /**
     * Extract and enhance code examples
     */
    function enhanceCodeExamples(html) {
        // Add copy buttons to code blocks
        const codeBlockRegex = /<pre class="language-([^"]*)"[^>]*><code[^>]*>([\s\S]*?)<\/code><\/pre>/gi;
        
        return html.replace(codeBlockRegex, (match, language, code) => {
            return `<div class="code-example">
                <div class="code-header">
                    <span class="code-language">${language}</span>
                    <button class="code-copy" onclick="copyCodeBlock(this)" title="Copy code">
                        <span class="copy-icon">📋</span>
                    </button>
                </div>
                ${match}
            </div>`;
        });
    }

    /**
     * Add table of contents links
     */
    function addTocLinks(html) {
        const headingRegex = /<h([2-4])[^>]*id="([^"]*)"[^>]*>(.*?)<\/h[2-4]>/gi;
        
        return html.replace(headingRegex, (match, level, id, text) => {
            return match.replace('>', ` data-toc-level="${level}" data-toc-id="${id}">`);
        });
    }

    /**
     * Sanitize HTML content (basic sanitization)
     */
    function sanitizeHtml(html) {
        // Remove potentially dangerous attributes and tags
        const dangerous = [
            /on\w+\s*=/gi,  // Remove event handlers
            /<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi,  // Remove scripts
            /<iframe\b[^<]*(?:(?!<\/iframe>)<[^<]*)*<\/iframe>/gi   // Remove iframes
        ];

        dangerous.forEach(pattern => {
            html = html.replace(pattern, '');
        });

        return html;
    }

    /**
     * Process ADR metadata
     */
    function extractMetadata(content) {
        const metadata = {
            title: null,
            status: null,
            lastModified: null,
            estimatedReadTime: null
        };

        // Extract title
        const titleMatch = content.match(/^#\s+(.+)$/m);
        if (titleMatch) {
            metadata.title = titleMatch[1].trim();
        }

        // Extract status
        const statusMatch = content.match(/##\s+Status\s*\n\s*([^\n]+)/i);
        if (statusMatch) {
            metadata.status = statusMatch[1].trim();
        }

        // Calculate estimated read time (average 200 words per minute)
        const wordCount = content.split(/\s+/).length;
        metadata.estimatedReadTime = Math.max(1, Math.round(wordCount / 200));

        return metadata;
    }

    /**
     * Add navigation improvements
     */
    function addNavigationEnhancements(html, adr) {
        // Add "back to top" links
        const sectionRegex = /(<h2[^>]*>.*?<\/h2>)/gi;
        
        html = html.replace(sectionRegex, (match) => {
            return `${match}<a href="#top" class="back-to-top" title="Back to top">↑</a>`;
        });

        // Add section navigation
        const sections = [];
        const sectionHeaderRegex = /<h2[^>]*id="([^"]*)"[^>]*>(.*?)<\/h2>/gi;
        let sectionMatch;
        
        while ((sectionMatch = sectionHeaderRegex.exec(html)) !== null) {
            sections.push({
                id: sectionMatch[1],
                title: sectionMatch[2].replace(/<[^>]*>/g, '')
            });
        }

        if (sections.length > 2) {
            const nav = `<nav class="adr-section-nav">
                <h4>Jump to Section</h4>
                <ul>
                    ${sections.map(section => 
                        `<li><a href="#${section.id}">${section.title}</a></li>`
                    ).join('')}
                </ul>
            </nav>`;
            
            // Insert after the title
            html = html.replace(/(<h1[^>]*>.*?<\/h1>)/i, `$1${nav}`);
        }

        return html;
    }

    /**
     * Format date for display
     */
    function formatDate(date) {
        if (!date) return null;
        return new Intl.DateTimeFormat('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        }).format(new Date(date));
    }

    /**
     * Copy code block to clipboard
     */
    window.copyCodeBlock = function(button) {
        const codeBlock = button.closest('.code-example').querySelector('code');
        const text = codeBlock.textContent;

        navigator.clipboard.writeText(text).then(() => {
            const icon = button.querySelector('.copy-icon');
            const originalText = icon.textContent;
            icon.textContent = '✓';
            setTimeout(() => {
                icon.textContent = originalText;
            }, 1500);
        }).catch(err => {
            console.error('Failed to copy code:', err);
        });
    };

    // Public API
    return {
        renderMarkdown,
        extractMetadata,
        configureMarked,
        sanitizeHtml,
        enhanceCodeExamples,
        formatDate
    };

})();