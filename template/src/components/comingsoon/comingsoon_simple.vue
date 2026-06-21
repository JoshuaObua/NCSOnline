<template>
	<div>
	
		<div class="page-wrapper">
		
			<div class="container-fluid p-0">
				<div class="comingsoon">
					<div class="comingsoon-inner text-center"><img src="../../assets/images/logo/logo-1.png" alt="">
						<h5>WE ARE COMING SOON</h5>
						<div class="countdown">
							              <Timer/>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
  import Timer from './comingsoon_timer';
  
  export default{
    components: {
      Timer
    },
    data(){
      return{
        startTime: 'Apr 29, 2021 00:00:00',
        endTime: 'Apr 29, 2022 00:00:00',
      };
    },
  };
</script>