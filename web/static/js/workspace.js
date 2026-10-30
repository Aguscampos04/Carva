// Presentation helpers. Business decisions remain in Go.
const CarvaWorkspace = (() => {
  const views = ['resumen', 'backlog', 'sprints', 'equipo', 'datos', 'reportes'];
  const tabs = ['info', 'poker', 'effort', 'defects'];
  function route(project, view = 'resumen', story = '', tab = 'info') {
    return project ? `#/proyecto/${encodeURIComponent(project)}/${view}${story ? `/${encodeURIComponent(story)}/${tab}` : ''}` : '#/';
  }
  function parseRoute(hash) {
    if (!hash || hash === '#' || hash === '#/') return {view:'inicio', project:''};
    try {
      const p = hash.replace(/^#\//, '').split('/').map(decodeURIComponent);
      if (p[0] !== 'proyecto' || !p[1] || !views.includes(p[2])) return null;
      if (p.length !== 3 && !(p.length === 5 && p[2] === 'backlog' && p[3] && tabs.includes(p[4]))) return null;
      return {project:p[1], view:p[2], story:p[3] || '', tab:p[4] || 'info'};
    } catch { return null; }
  }
  const normalize = text => String(text).normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase().trim();
  function filterStories(items, search = '', priority = '', status = '') {
    return items.filter(item => normalize(item.title).includes(normalize(search)) && (!priority || item.priority === priority) && (!status || item.status === status));
  }
  function createGate() {
    let generation = 0;
    return {next:() => ++generation, current:value => value === generation};
  }
  function points(item) {
    if (item.story_points != null) return `${item.story_points} SP acordados`;
    if (item.provisional_story_points != null) return `${item.provisional_story_points} SP provisionales`;
    return 'Sin estimar';
  }
  function pendingFirst(items) {
    const completed = item => item.acceptance_criteria_verified === true || item.status === 'Completada';
    return [...items.filter(item => !completed(item)), ...items.filter(completed)];
  }
  return {route, parseRoute, filterStories, createGate, points, pendingFirst};
})();
