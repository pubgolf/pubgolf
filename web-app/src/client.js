// Import global/unscoped CSS files these here so they get included in the main
// bundle.
import 'src/assets/css/tailwind.css';
import 'src/assets/css/global.css';

import { start } from '@sapper/app';

start({
  target: document.querySelector('#sapper'),
});
