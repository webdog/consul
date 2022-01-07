(routes => routes({
  dc: {
    nspaces: {
      _options: {
        path: '/namespaces',
        abilities: ['read nspaces'],
      },
      index: {
        _options: {
          path: '/',
          queryParams: {
            sortBy: 'sort',
            searchproperty: {
              as: 'searchproperty',
              empty: [['Name', 'Description', 'Role', 'Policy']],
            },
            search: {
              as: 'filter',
              replace: true,
            },
          },
        },
      },
      edit: {
        _options: { path: '/:name' },
      },
      create: {
        _options: {
          template: 'dc/nspaces/edit',
          path: '/create',
          abilities: ['create nspaces'],
        },
      },
    },
  },
}))(
  (json, data = (typeof document !== 'undefined' ? document.currentScript.dataset : module.exports)) => {
    data[`routes`] = JSON.stringify(json);
  }
);
