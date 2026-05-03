import React from "react";

const SEO_SCORE_CONFIG = {
  high: { color: "text-teal-400",   bg: "bg-teal-500/20 border-teal-500/40",    label: "Bueno" },
  mid:  { color: "text-orange-400", bg: "bg-orange-500/20 border-orange-500/40", label: "Mejorable" },
  low:  { color: "text-red-400",    bg: "bg-red-500/20 border-red-500/40",       label: "Deficiente" },
};

function getSEOConfig(score) {
  if (score >= 70) return SEO_SCORE_CONFIG.high;
  if (score >= 40) return SEO_SCORE_CONFIG.mid;
  return SEO_SCORE_CONFIG.low;
}

function Badge({ children, color = "gray" }) {
  const colors = {
    teal:   "bg-teal-500/20 text-teal-400",
    orange: "bg-orange-500/20 text-orange-400",
    red:    "bg-red-500/20 text-red-400",
    gray:   "bg-gray-500/20 text-gray-400",
    blue:   "bg-blue-500/20 text-blue-400",
  };
  return (
    <span className={`text-xs px-1.5 py-0.5 rounded font-medium ${colors[color] ?? colors.gray}`}>
      {children}
    </span>
  );
}

function Section({ title, children }) {
  return (
    <div className="bg-white/5 rounded-lg p-4 mb-6">
      <h4 className="text-lg font-semibold text-white mb-3">{title}</h4>
      {children}
    </div>
  );
}

function KV({ label, children }) {
  return (
    <div>
      <span className="text-gray-400">{label}:</span>
      <div className="text-white mt-0.5">{children}</div>
    </div>
  );
}

export default function DetailModal({ result, onClose }) {
  if (!result) return null;

  const {
    url, status_code, content_type, load_time_ms, word_count, created_at,
    title, description, keywords, author, language, site_name,
    headers, links, images, content,
    canonical_url, robots_directive, x_robots_tag, viewport,
    og_data, twitter_card, schema_org, redirect_chain, final_url,
    h1_count, has_multiple_h1, seo_score,
  } = result;

  const safeHeaders   = headers        ?? [];
  const safeLinks     = links          ?? [];
  const safeImages    = images         ?? [];
  const safeContent   = content        ?? "";
  const safeSchemaOrg = schema_org     ?? [];
  const safeRedirects = redirect_chain ?? [];
  const safeOG        = og_data        ?? {};
  const safeTwitter   = twitter_card   ?? {};

  const date        = new Date(created_at).toLocaleString("es-ES");
  const statusColor = status_code === 200 ? "text-teal-400" : "text-red-400";
  const seoConfig   = getSEOConfig(seo_score ?? 0);

  const isNoindex = [robots_directive, x_robots_tag]
    .some((v) => v?.toLowerCase().includes("noindex"));

  const internalLinks      = safeLinks.filter((l) => l.is_internal);
  const externalLinks      = safeLinks.filter((l) => !l.is_internal);
  const imagesWithoutAlt   = safeImages.filter(
    (img) => !img.alt && !(img.src ?? img)?.startsWith("data:")
  );

  return (
    <div
      id="detailModal"
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
      onClick={(e) => e.target.id === "detailModal" && onClose()}
    >
      <div className="animate-modal-in bg-black/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 max-w-4xl w-full max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between px-6 py-5 border-b border-white/20 flex-shrink-0">
          <h3 className="text-2xl font-bold text-white">Detalles del Scraping</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors p-1"
            aria-label="Cerrar"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="overflow-y-auto p-6">

          {/* SEO Score */}
          {seo_score !== undefined && (
            <div className={`rounded-lg p-4 mb-6 border ${seoConfig.bg} flex items-center gap-4`}>
              <div className="text-center min-w-[56px]">
                <div className={`text-4xl font-bold ${seoConfig.color}`}>{seo_score}</div>
                <div className={`text-xs font-medium ${seoConfig.color}`}>/100</div>
              </div>
              <div>
                <div className={`text-lg font-semibold ${seoConfig.color}`}>
                  SEO Score — {seoConfig.label}
                </div>
                <div className="flex flex-wrap gap-1.5 mt-2">
                  {!canonical_url && <Badge color="orange">Sin canonical</Badge>}
                  {isNoindex && <Badge color="red">noindex</Badge>}
                  {h1_count === 0 && <Badge color="red">Sin H1</Badge>}
                  {has_multiple_h1 && <Badge color="orange">Múltiples H1</Badge>}
                  {safeSchemaOrg.length === 0 && <Badge color="gray">Sin JSON-LD</Badge>}
                  {imagesWithoutAlt.length > 0 && (
                    <Badge color="orange">{imagesWithoutAlt.length} img sin alt</Badge>
                  )}
                  {safeRedirects.length > 0 && (
                    <Badge color="orange">
                      {safeRedirects.length} redirect{safeRedirects.length > 1 ? "s" : ""}
                    </Badge>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Información Básica */}
          <Section title="Información Básica">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-300">
              <div>
                <span className="text-gray-400">URL:</span>
                <p className="text-cyan-400 break-all">{url}</p>
              </div>
              {final_url && final_url !== url && (
                <div>
                  <span className="text-gray-400">URL final (post-redirect):</span>
                  <p className="text-cyan-300 break-all">{final_url}</p>
                </div>
              )}
              <div>
                <span className="text-gray-400">Código de estado:</span>
                <p className={`${statusColor} font-medium`}>{status_code}</p>
              </div>
              <div>
                <span className="text-gray-400">Tipo de contenido:</span>
                <p className="text-white">{content_type || "Desconocido"}</p>
              </div>
              <div>
                <span className="text-gray-400">Tiempo de carga:</span>
                <p className="text-white">{load_time_ms}ms</p>
              </div>
              <div>
                <span className="text-gray-400">Número de palabras:</span>
                <p className="text-white">{word_count || 0}</p>
              </div>
              <div>
                <span className="text-gray-400">Scraped:</span>
                <p className="text-white">{date}</p>
              </div>
            </div>
          </Section>

          {/* Indexabilidad */}
          <Section title="Indexabilidad">
            <div className="space-y-2 text-sm">
              <div className="flex items-start gap-2">
                <span className="text-gray-400 w-32 flex-shrink-0">Canonical:</span>
                {canonical_url ? (
                  <a href={canonical_url} target="_blank" rel="noopener noreferrer"
                    className="text-cyan-400 hover:text-cyan-300 break-all">
                    {canonical_url}
                  </a>
                ) : (
                  <span className="text-orange-400">No definida</span>
                )}
              </div>
              <div className="flex items-center gap-2">
                <span className="text-gray-400 w-32 flex-shrink-0">Meta robots:</span>
                {robots_directive ? (
                  <span className={isNoindex ? "text-red-400 font-medium" : "text-teal-400"}>
                    {robots_directive}
                  </span>
                ) : (
                  <span className="text-gray-500">No definido (indexable por defecto)</span>
                )}
              </div>
              {x_robots_tag && (
                <div className="flex items-center gap-2">
                  <span className="text-gray-400 w-32 flex-shrink-0">X-Robots-Tag:</span>
                  <span className={x_robots_tag.toLowerCase().includes("noindex") ? "text-red-400" : "text-white"}>
                    {x_robots_tag}
                  </span>
                </div>
              )}
              <div className="flex items-center gap-2">
                <span className="text-gray-400 w-32 flex-shrink-0">Viewport:</span>
                {viewport ? (
                  <span className="text-teal-400 break-all">{viewport}</span>
                ) : (
                  <span className="text-orange-400">No definido (posible problema mobile)</span>
                )}
              </div>
              <div className="flex items-center gap-2">
                <span className="text-gray-400 w-32 flex-shrink-0">H1:</span>
                {h1_count === 0 ? (
                  <Badge color="red">Ningún H1 encontrado</Badge>
                ) : has_multiple_h1 ? (
                  <Badge color="orange">{h1_count} H1 (debería ser exactamente 1)</Badge>
                ) : (
                  <Badge color="teal">1 H1 ✓</Badge>
                )}
              </div>
            </div>
          </Section>

          {/* Metadatos */}
          <Section title="Metadatos">
            <div className="space-y-2 text-sm text-gray-300">
              <KV label="Título">
                <span>{title || "Sin título"}</span>
                {title && (
                  <span className={`ml-2 text-xs ${
                    title.length >= 50 && title.length <= 60 ? "text-teal-400" : "text-orange-400"
                  }`}>
                    ({title.length} chars — ideal 50-60)
                  </span>
                )}
              </KV>
              <KV label="Descripción">
                <span>{description || "Sin descripción"}</span>
                {description && (
                  <span className={`ml-2 text-xs ${
                    description.length >= 150 && description.length <= 160 ? "text-teal-400" : "text-orange-400"
                  }`}>
                    ({description.length} chars — ideal 150-160)
                  </span>
                )}
              </KV>
              <KV label="Palabras clave">{keywords || "Sin palabras clave"}</KV>
              <KV label="Autor">{author || "Desconocido"}</KV>
              <KV label="Idioma">{language || "Desconocido"}</KV>
              <KV label="Nombre del sitio">{site_name || "Desconocido"}</KV>
            </div>
          </Section>

          {/* Open Graph */}
          {(safeOG.title || safeOG.description || safeOG.url || safeOG.type) && (
            <Section title="Open Graph">
              <div className="space-y-2 text-sm text-gray-300">
                {safeOG.title && <KV label="og:title">{safeOG.title}</KV>}
                {safeOG.description && <KV label="og:description">{safeOG.description}</KV>}
                {safeOG.url && (
                  <KV label="og:url">
                    <a href={safeOG.url} target="_blank" rel="noopener noreferrer"
                      className="text-cyan-400 hover:text-cyan-300 break-all">
                      {safeOG.url}
                    </a>
                  </KV>
                )}
                {safeOG.type && <KV label="og:type">{safeOG.type}</KV>}
                {safeOG.image && (
                  <KV label="og:image">
                    <a href={safeOG.image} target="_blank" rel="noopener noreferrer"
                      className="text-cyan-400 hover:text-cyan-300 break-all">
                      {safeOG.image}
                    </a>
                  </KV>
                )}
                {safeOG.site_name && <KV label="og:site_name">{safeOG.site_name}</KV>}
                {safeOG.locale && <KV label="og:locale">{safeOG.locale}</KV>}
              </div>
            </Section>
          )}

          {/* Twitter Card */}
          {(safeTwitter.card || safeTwitter.title || safeTwitter.description) && (
            <Section title="Twitter Card">
              <div className="space-y-2 text-sm text-gray-300">
                {safeTwitter.card && <KV label="twitter:card">{safeTwitter.card}</KV>}
                {safeTwitter.title && <KV label="twitter:title">{safeTwitter.title}</KV>}
                {safeTwitter.description && <KV label="twitter:description">{safeTwitter.description}</KV>}
                {safeTwitter.image && (
                  <KV label="twitter:image">
                    <a href={safeTwitter.image} target="_blank" rel="noopener noreferrer"
                      className="text-cyan-400 hover:text-cyan-300 break-all">
                      {safeTwitter.image}
                    </a>
                  </KV>
                )}
                {safeTwitter.site && <KV label="twitter:site">{safeTwitter.site}</KV>}
              </div>
            </Section>
          )}

          {/* JSON-LD */}
          {safeSchemaOrg.length > 0 && (
            <Section title={`Datos Estructurados — JSON-LD (${safeSchemaOrg.length})`}>
              <div className="space-y-3">
                {safeSchemaOrg.map((raw, i) => {
                  let type = "Desconocido";
                  try {
                    const parsed = JSON.parse(raw);
                    type = parsed["@type"] ?? "Desconocido";
                  } catch { /* JSON inválido — se muestra tipo "Desconocido" */ }
                  return (
                    <div key={i} className="text-sm">
                      <div className="mb-1">
                        <Badge color="blue">
                          @type: {Array.isArray(type) ? type.join(", ") : String(type)}
                        </Badge>
                      </div>
                      <pre className="text-gray-400 text-xs whitespace-pre-wrap break-all max-h-32 overflow-y-auto bg-black/30 rounded p-2">
                        {raw.length > 500 ? raw.substring(0, 500) + "…" : raw}
                      </pre>
                    </div>
                  );
                })}
              </div>
            </Section>
          )}

          {/* Cadena de redirects */}
          {safeRedirects.length > 0 && (
            <Section title={`Cadena de Redirects (${safeRedirects.length})`}>
              <div className="space-y-1 text-sm font-mono">
                <div className="text-cyan-400 break-all">{url}</div>
                {safeRedirects.map((u, i) => (
                  <div key={i} className="flex items-center gap-2">
                    <span className="text-orange-400">→</span>
                    <a href={u} target="_blank" rel="noopener noreferrer"
                      className="text-cyan-300 hover:text-cyan-200 break-all text-xs">
                      {u}
                    </a>
                  </div>
                ))}
              </div>
            </Section>
          )}

          {/* Cabeceras HTML */}
          {safeHeaders.length > 0 && (
            <Section title={`Cabeceras HTML (${safeHeaders.length})`}>
              <div className="max-h-40 overflow-y-auto text-sm space-y-1">
                {safeHeaders.map((h, i) => (
                  <div key={i} style={{ paddingLeft: `${(h.level - 1) * 12}px` }}>
                    <span className="text-teal-400 font-mono text-xs">H{h.level}</span>{" "}
                    <span className="text-white">{h.text}</span>
                  </div>
                ))}
              </div>
            </Section>
          )}

          {/* Links */}
          {safeLinks.length > 0 && (
            <Section
              title={`Links (${safeLinks.length} — ${internalLinks.length} internos, ${externalLinks.length} externos)`}
            >
              <div className="max-h-48 overflow-y-auto text-sm space-y-2">
                {safeLinks.map((link, i) => {
                  const href = link.url ?? link;
                  const anchor = link.anchor_text;
                  const isInt = link.is_internal;
                  const isNofollow = link.rel?.includes("nofollow");
                  return (
                    <div key={i} className="flex items-start gap-1.5">
                      <Badge color={isInt ? "teal" : "gray"}>{isInt ? "INT" : "EXT"}</Badge>
                      {isNofollow && <Badge color="orange">NF</Badge>}
                      <div className="flex-1 min-w-0">
                        <a href={href} target="_blank" rel="noopener noreferrer"
                          className="text-cyan-400 hover:text-cyan-300 truncate block">
                          {href}
                        </a>
                        {anchor && (
                          <span className="text-gray-500 text-xs">"{anchor}"</span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            </Section>
          )}

          {/* Imágenes */}
          {safeImages.length > 0 && (
            <Section
              title={`Imágenes (${safeImages.length}${imagesWithoutAlt.length > 0 ? ` — ${imagesWithoutAlt.length} sin alt` : ""})`}
            >
              <div className="max-h-48 overflow-y-auto text-sm space-y-2">
                {safeImages.map((img, i) => {
                  const src = img.src ?? img;
                  const alt = img.alt;
                  return (
                    <div key={i} className="flex items-start gap-1.5">
                      {!alt && <Badge color="red">sin alt</Badge>}
                      <div className="flex-1 min-w-0">
                        <a href={src} target="_blank" rel="noopener noreferrer"
                          className="text-teal-400 hover:text-teal-300 truncate block">
                          {src}
                        </a>
                        {alt && <span className="text-gray-500 text-xs">Alt: "{alt}"</span>}
                        {img.title && (
                          <span className="text-gray-600 text-xs ml-2">Title: "{img.title}"</span>
                        )}
                      </div>
                    </div>
                  );
                })}
              </div>
            </Section>
          )}

          {/* Vista previa del contenido */}
          {safeContent && (
            <Section title="Vista previa del contenido">
              <pre className="whitespace-pre-wrap break-words text-sm text-gray-300 max-h-60 overflow-y-auto">
                {safeContent.length > 2000
                  ? `${safeContent.substring(0, 2000)}...\n\n[Content truncated]`
                  : safeContent}
              </pre>
            </Section>
          )}

        </div>
      </div>
    </div>
  );
}
