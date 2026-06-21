<template>
    <h6 class="mb-0">Personal Details</h6>
										<form class="form theme-form">
											<div class="row">
												<div class="col">
													<div class="form-group">
														<label for="exampleFormControlInput1" class="form-label">Full Name:<span class="font-danger">*</span></label>
														<input class="form-control" id="exampleFormControlInput1" type="text" placeholder="Enter your full name">
													</div>
												</div>
											</div>
											<div class="row">
												<div class="col">
													<div class="form-group">
														<label for="exampleFormControlInput2" class="form-label">Email:<span class="font-danger">*</span></label>
														<input class="form-control" id="exampleFormControlInput2" type="email" placeholder="Enter email">
													</div>
												</div>
											</div>
											<div class="row">
												<div class="col">
													<div class="form-group">
														<label for="exampleFormControlInput3" class="form-label">Password:<span class="font-danger">*</span></label>
														<input class="form-control" id="exampleFormControlInput3" type="password" placeholder="Enter password">
													</div>
												</div>
											</div>
											<div class="row">
												<div class="col">
													<div class="form-group">
														<label for="exampleFormControlpassword" class="form-label">Repeat Password:<span class="font-danger">*</span>
														</label>
														<input class="form-control" id="exampleFormControlpassword" type="password" placeholder="Repeat password">
													</div>
												</div>
											</div>
											<div class="row birthday-select">
												<birthdayselect/>
											</div>
										</form>
</template>
<script>
import birthdayselect from "./birthdayselect.vue"
export default {
	components:{
      birthdayselect
	}
}
</script>