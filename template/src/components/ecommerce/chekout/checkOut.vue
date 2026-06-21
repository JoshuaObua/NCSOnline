<template>
  <div class="checkout">
    <Breadcrumbs main="Ecommerce" title="Checkout" />
    <div class="container-fluid ">
      <div class="card">
        <div class="card-header pb-0">
          <h5>Billing Details</h5>
        </div>
        <div class="card-body">
          <form>
          <div class="row">
            <div class="col-lg-6 col-sm-12">
              <div class="row">
                <div class="form-group col-sm-6">
                  <label for="inputEmail4">First Name</label>
                  <input type="text" v-model="user.firstName" id="firstName" name="firstName" class="form-control" :class="{ 'is-invalid': submitted && errors.has('firstName') }" />
                  <div v-if="submitted && errors.has('firstName')" class="invalid-feedback"
                  >{{ errors.first('firstName') }}</div>
                </div>
                <div class="form-group col-sm-6">
                  <label for="inputPassword4">Last Name</label>
                  <input type="text" v-model="user.lastName" id="lastName" name="lastName" class="form-control" :class="{ 'is-invalid': submitted && errors.has('lastName') }" />
                </div>
              </div>
              <div class="row">
                <div class="form-group col-sm-6">
                  <label for="inputEmail4">Phone</label>
                  <input type="text"  v-model="user.phone" id="phone" name="phone" class="form-control" :class="{ 'is-invalid': submitted && errors.has('phone') }" />
                </div>
                <div class="form-group col-sm-6">
                  <label for="inputPassword4">Email Address</label>
                  <input type="email" v-model="user.email" id="email" name="email" class="form-control" :class="{ 'is-invalid': submitted && errors.has('email') }"  />
                </div>
              </div>
              <div class="form-group">
                <label for="inputState">Country</label>
                <select class="form-control form-select" id="inputState">
                  <option selected>Choose...</option>
                  <option >India</option>
                  <option>USA</option>
                  <option>UK</option>
                </select>
              </div>
              <div class="form-group">
                <label for="inputAddress2">Address</label>
                <input type="text"  v-model="user.address"  id="address" name="address" class="form-control" :class="{ 'is-invalid': submitted && errors.has('address') }" />
              </div>
              <div class="form-group">
                <label for="inputCity">Town/City</label>
                <input type="text" v-model="user.city"  id="city" name="city" class="form-control" :class="{ 'is-invalid': submitted && errors.has('city') }" />
              </div>
              <div class="form-group">
                <label for="inputAddress2">State/Country</label>
                <input type="text" v-model="user.state" id="state" name="state" class="form-control" :class="{ 'is-invalid': submitted && errors.has('state') }" />
              </div>
              <div class="form-group">
                <label for="inputAddress2">Postal Code</label>
                <input type="text" v-model="user.pincode"  id="pincode" name="pincode" class="form-control" :class="{ 'is-invalid': submitted && errors.has('pincode') }"/>
              </div>
               <div class="mb-3">
                        <div class="form-check">
                          <input class="form-check-input" id="gridCheck" type="checkbox">
                          <label class="form-check-label" for="gridCheck">Check me out</label>
                        </div>
                      </div>
            </div>
            <checkoutDetails/>
          </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import checkoutDetails from "./checkoutDetails.vue"
  export default {
    components: {
     checkoutDetails
    },
    data() {
      return {
        user: {
          firstName: '',
          lastName: '',
          phone: '',
          email: '',
          address: '',
          city: '',
          state: '',
          pincode: ''
        },
      };
    },
  };
</script>
