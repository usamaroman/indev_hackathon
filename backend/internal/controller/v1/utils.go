package v1

// func buildFilters(c *gin.Context) *entity.CarFilters {
// 	businessID := businessID(c)
// 	limit, errLimit := limit(c)
// 	offset, errOffset := offset(c)
// 	yearFromStr := c.Query("year_from")
// 	yearToStr := c.Query("year_to")
// 	priceFromStr := c.Query("price_from")
// 	priceToStr := c.Query("price_to")
// 	sortPrice := c.Query("sort")
// 	brandStr := c.Query("brand")
// 	classStr := c.Query("class")
// 	transmissionStr := c.Query("transmission")
// 	startDateStr := c.Query("start_date")
// 	endDateStr := c.Query("end_date")
// 	role := c.GetString(middleware.UserTypeCtx)
// 	queryUrl := c.Request.URL.RawQuery
//
// 	ok := checkLimitOffset(limit, offset)
// 	if !ok || errLimit != nil || errOffset != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": fmt.Sprintf("Bad limit: %d, offest: %d", limit, offset),
// 		})
// 		return nil
// 	}
//
// 	var startDate, endDate time.Time
// 	var err error
//
// 	if startDateStr != "" {
// 		startDate, err = time.Parse("02.01.2006", startDateStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return nil
// 		}
// 	}
//
// 	if endDateStr != "" {
// 		endDate, err = time.Parse("02.01.2006", endDateStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": err.Error(),
// 			})
// 			return nil
// 		}
// 	}
//
// 	var brands []types.CarBrand
// 	var classes []types.CarClass
// 	var transmissions []types.CarTransmission
// 	var priceFrom float32
// 	var priceTo float32
//
// 	if brandStr != "" {
// 		brandsArray := strings.Split(brandStr, ",")
//
// 		for _, brand := range brandsArray {
// 			val, ok := types.StringToCarBrand[brand]
// 			if !ok {
// 				continue
// 			}
//
// 			brands = append(brands, val)
// 		}
// 	}
//
// 	if classStr != "" {
// 		classArray := strings.Split(classStr, ",")
//
// 		for _, class := range classArray {
// 			val, ok := types.StringToCarClass[class]
// 			if !ok {
// 				continue
// 			}
//
// 			classes = append(classes, val)
// 		}
// 	}
//
// 	if transmissionStr != "" {
// 		transmissionArray := strings.Split(transmissionStr, ",")
//
// 		for _, transmission := range transmissionArray {
// 			val, ok := types.StringToCarTransmission[transmission]
// 			if !ok {
// 				continue
// 			}
//
// 			transmissions = append(transmissions, val)
// 		}
// 	}
//
// 	if priceFromStr != "" {
// 		priceFrom64, err := strconv.ParseFloat(priceFromStr, 32)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "failed to parse filter price from",
// 			})
// 			return nil
// 		}
// 		priceFrom = float32(priceFrom64)
// 	}
//
// 	if priceToStr != "" {
// 		priceTo64, err := strconv.ParseFloat(priceToStr, 32)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "failed to parse filter price to",
// 			})
// 			return nil
// 		}
// 		priceTo = float32(priceTo64)
// 	}
//
// 	if yearFromStr != "" {
// 		_, err := strconv.Atoi(yearFromStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "failed to parse filter year from",
// 			})
// 			return nil
// 		}
// 	}
//
// 	if yearToStr != "" {
// 		_, err := strconv.Atoi(yearToStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "failed to parse filter year to",
// 			})
// 			return nil
// 		}
// 	}
//
// 	filters := entity.CarFilters{
// 		Limit:        uint64(limit),
// 		Offset:       uint64(offset),
// 		BusinessID:   businessID,
// 		YearFrom:     yearFromStr,
// 		YearTo:       yearToStr,
// 		PriceFrom:    priceFrom,
// 		PriceTo:      priceTo,
// 		Brand:        brands,
// 		Class:        classes,
// 		Transmission: transmissions,
// 		StartDate:    startDate,
// 		EndDate:      endDate,
// 		Role:         role,
// 		QueryUrl:     queryUrl,
// 	}
//
// 	switch sortPrice {
// 	case "prc.d":
// 		filters.PriceDesc = true
// 	case "prc.a":
// 		filters.PriceAsc = true
// 	}
//
// 	return &filters
// }
