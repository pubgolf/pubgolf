<script>
  import { stores } from '@sapper/app';
  import { initClient } from '../../api';
  import { event } from '../../stores';


  const { page, session } = stores();

  const { params } = $page;
  initClient(params.event, $session);
  $event = params.event;
</script>

<slot></slot>
