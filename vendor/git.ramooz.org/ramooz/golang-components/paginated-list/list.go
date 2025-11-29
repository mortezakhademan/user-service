package componentsList

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"slices"
	"strings"

	"git.ramooz.org/ramooz/golang-components/paginated-list/mongodb/models"
)

const (
	DEFAULT_PAGE_SIZE          = 20
	RESULT_TYPE_PAGINATED_JSON = 0
	RESULT_TYPE_EXCEL          = 1
)

type pipeline struct {
	Pipeline              bson.D
	ChangeableResultCount bool // ChangeableResultCount this field is true when stage is match, group or another stages that can change result count
	ResultColumns         map[string]*models.ColumnInfo
}

// /use global.FormDecoder.Decode(list,ctx.FormValues()) to fill list
type List struct {
	Filters                  map[string]string             `form:"filters" json:"filters"`
	Sort                     string                        `form:"sort" json:"sort"`
	PageNo                   int                           `form:"page" json:"page"`
	PageSize                 int                           `form:"page-size" json:"-"`
	ResultType               int                           `form:"result-type" json:"-"`
	Data                     interface{}                   `form:"-" json:"data"`
	ExtraData                interface{}                   `form:"-" json:"extraData"`
	TotalItemsCount          int                           `form:"-" json:"totalItemsCount"`
	ExcelFile                *excelize.File                `form:"-" json:"-"`
	fromCollection           *mongo.Collection             `form:"-"`
	fromCollectionColumnsMap map[string]*models.ColumnInfo `form:"-"`
	exportColumnsInfo        []*models.ExportColumnInfo
	CurrentUserId            null.Int `form:"-" json:"-"`
	dbPipelines              []pipeline
	sortStageAdded           bool
	sortFields               []*models.SortField
}

func NewList() *List {
	return &List{
		PageSize: DEFAULT_PAGE_SIZE,
		PageNo:   1,
	}
}

func (list *List) SetExportColumns(columns []*models.ExportColumnInfo) {
	list.exportColumnsInfo = columns
}

// AddPipelineStages add multiple stages to pipeline and ChangeableResultCount this field is true when stage is match, group or another filter stages
func (list *List) AddPipelineStages(stages []bson.D, resultColumns map[string]*models.ColumnInfo, hasChangeableResultCount bool) {
	if len(stages) > 1 {
		for _, stage := range stages[:len(stages)-1] {
			list.AddPipelineStage(stage, nil, false)
		}
	}
	/// add resultColumns after all pipelines
	if len(stages) > 0 {
		list.dbPipelines = append(list.dbPipelines, pipeline{
			Pipeline:              stages[len(stages)-1],
			ResultColumns:         resultColumns,
			ChangeableResultCount: hasChangeableResultCount,
		})
	}
}

func (list *List) AddMultilingualFilter(multilingualFields []string) {
	//aggregate := map[string]bson.D{}
	//
	//for _, field := range multilingualFields { // lang-title
	//	fieldData := strings.Split(field, "-")
	//	if value, ok := list.Filters[fieldData[0]+"-"+fieldData[1]]; ok {
	//		if _,ok := aggregate[fieldData[1]];!ok {
	//			aggregate[fieldData[1]] = bson.D{}
	//		}
	//		aggregate[fieldData[1]] = bson.D{
	//			{
	//				"$and", bson.D{
	//				{
	//					fieldData[1] + ".lang", fieldData[0],
	//				}, {
	//					fieldData[1] + ".value", value,
	//				},
	//			},
	//			},
	//		}
	//	}
	//}
	//list.dbPipelines = append(list.dbPipelines, aggregate...)
}

// AddPipelineStage add multiple stages to pipeline and ChangeableResultCount this field is true when stage is match, group or another filter stages
func (list *List) AddPipelineStage(stage bson.D, resultColumns map[string]*models.ColumnInfo, isChangeableResultCount bool) {
	list.dbPipelines = append(list.dbPipelines, pipeline{
		Pipeline:              stage,
		ResultColumns:         resultColumns,
		ChangeableResultCount: isChangeableResultCount,
	})
}

func (list *List) prepareSortFields() {
	if list.Sort == "" {
		return
	}
	sortFields := strings.Split(list.Sort, ",")
	for _, sortField := range sortFields {
		field, direction := parseSortField(sortField)
		list.sortFields = append(list.sortFields, &models.SortField{
			Column:    field,
			Direction: direction,
		})
	}
}
func (list *List) RunQuery(ctx context.Context, fromCollection *mongo.Collection, fromCollectionColumnsMap map[string]*models.ColumnInfo, defaultSort string, dataSliceType interface{}) {
	list.fromCollection = fromCollection
	list.fromCollectionColumnsMap = fromCollectionColumnsMap
	if list.Sort == "" {
		list.Sort = defaultSort
	}
	list.prepareSortFields()
	dbPipelines, countPipelines := list.mergePaginationPipelines()
	list.setListTotalCount(countPipelines, ctx)

	cursor, err := fromCollection.Aggregate(ctx, dbPipelines)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(ctx)
	list.prepareResult(cursor, ctx, dataSliceType)
}

func (list *List) generateFilterStage(filters map[string]string, columnsMap map[string]*models.ColumnInfo) (stage bson.D, remainingFilters map[string]string) {
	filterStage := bson.D{}
	for filterCol, filterValue := range filters {
		if columnInfo, ok := columnsMap[filterCol]; ok {
			filterStage = append(filterStage, columnInfo.GetFilterVal(filterValue))
			delete(filters, filterCol)
		}
	}
	if len(filterStage) > 0 {
		stage = bson.D{{"$match", filterStage}}
	}
	return stage, filters
}

func (list *List) mergePaginationPipelines() (dataPipeline mongo.Pipeline, countPipeline mongo.Pipeline) {
	allColumnsMap := list.fromCollectionColumnsMap

	preparedPipelines := mongo.Pipeline{}
	remainingFilters := make(map[string]string)
	for k, v := range list.Filters {
		remainingFilters[k] = v
	}
	/// add main collection filters
	filterStage, remainingFilters := list.generateFilterStage(remainingFilters, list.fromCollectionColumnsMap)
	if len(filterStage) > 0 {
		preparedPipelines = append(preparedPipelines, filterStage)
	}

	/// add pipelines before sort and pagination stages that maybe change result count
	stages, pipelineIndex, remainingFilters := list.addChangeableResultCountStagesAndFilters(remainingFilters)
	preparedPipelines = append(preparedPipelines, stages...)
	/// add pipelines until remaining filters make empty because pagination stages must add after these pipelines
	for len(remainingFilters) != 0 && pipelineIndex < len(list.dbPipelines) {
		currentStage := list.dbPipelines[pipelineIndex]
		preparedPipelines = append(preparedPipelines, currentStage.Pipeline)
		pipelineIndex++
		if len(currentStage.ResultColumns) > 0 {
			for k, v := range currentStage.ResultColumns {
				allColumnsMap[k] = v
			}
			filterStage, remainingFilters = list.generateFilterStage(remainingFilters, currentStage.ResultColumns)
			if len(filterStage) > 0 {
				preparedPipelines = append(preparedPipelines, filterStage)
			}
		}
	}

	countPipeline = append(slices.Clone(preparedPipelines), list.getCountStage())
	// add pagination stages after others stages that can change result count
	preparedPipelines = list.addSortAndPageLimitStages(preparedPipelines, allColumnsMap)

	/// add remaining pipelines
	for i := pipelineIndex; i < len(list.dbPipelines); i++ {
		pipeline := list.dbPipelines[i]
		preparedPipelines = append(preparedPipelines, pipeline.Pipeline)
		if pipeline.ResultColumns != nil {
			for k, v := range pipeline.ResultColumns {
				allColumnsMap[k] = v
			}
			preparedPipelines = list.addSortAndPageLimitStages(preparedPipelines, allColumnsMap)
		}
	}

	/// if client send incorrect sort field pagination stage must to add pipelines
	if !list.sortStageAdded && list.isPaginatedList() {
		log.Printf("sort %s is wrong for sort stage, collection name: %s", list.Sort, list.fromCollection.Name())
		preparedPipelines = append(preparedPipelines, list.getPageLimitStages()...)
	}
	return preparedPipelines, countPipeline
}

func pipelinesHasChangeableResultCount(pipelines []pipeline) bool {
	for _, p := range pipelines {
		if p.ChangeableResultCount {
			return true
		}
	}
	return false
}

func (list *List) addSortAndPageLimitStages(preparedPipelines mongo.Pipeline, columnsMap map[string]*models.ColumnInfo) mongo.Pipeline {
	if !list.sortStageAdded {
		if list.Sort == "" {
			if list.isPaginatedList() {
				preparedPipelines = append(preparedPipelines, list.getPageLimitStages()...)
			}
			list.sortStageAdded = true
			return preparedPipelines
		}
		if sortStage := list.getSortStage(columnsMap); sortStage != nil {
			list.sortStageAdded = true
			preparedPipelines = append(preparedPipelines, sortStage)
			if list.isPaginatedList() {
				preparedPipelines = append(preparedPipelines, list.getPageLimitStages()...)
			}
		}
	}
	return preparedPipelines
}

func (list *List) getFindOptions(columnsMap map[string]models.ColumnInfo) *options.FindOptionsBuilder {

	findOptions := options.Find()
	findOptions.SetLimit(int64(list.PageSize))
	findOptions.SetSkip(int64((list.PageNo - 1) * list.PageSize))
	if list.Sort != "" {
		sortDirection := 1
		sort := list.Sort
		if strings.Index(list.Sort, "-") == 0 {
			sort = list.Sort[1:]
			sortDirection = -1
		}
		findOptions.SetSort(bson.D{{sort, sortDirection}})
	}
	return findOptions
}

func (list *List) getAggregateOptions(columnsMap map[string]models.ColumnInfo) *options.AggregateOptionsBuilder {

	aggregateOptions := options.Aggregate()
	//aggregateOptions.SetLimit(int64(list.PageSize))
	//aggregateOptions.SetSkip(int64((list.PageNo - 1) * list.PageSize))
	//if list.Sort != "" {
	//	sortDirection := 1
	//	sort := list.Sort
	//	if strings.Index(list.Sort, "-") == 0 {
	//		sort = list.Sort[1:]
	//		sortDirection = -1
	//	}
	//	aggregateOptions.SetSort(bson.D{{sort, sortDirection}})
	//}
	return aggregateOptions
}

func (list *List) getCountStage() bson.D {
	return bson.D{{"$count", "count"}}
}

func (list *List) setListTotalCount(countPipeline mongo.Pipeline, context context.Context) {
	countCursor, err := list.fromCollection.Aggregate(context, countPipeline)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer countCursor.Close(context)
	countCursor.Next(context)
	countData := bson.M{}
	countCursor.Decode(&countData)
	if count, ok := countData["count"]; ok {
		list.TotalItemsCount = int(count.(int32))
	}
}

func (list *List) getPageLimitStages() []bson.D {
	if list.PageNo == 1 {
		return []bson.D{{{"$limit", list.PageSize}}}
	}
	return []bson.D{{{"$skip", (list.PageNo - 1) * list.PageSize}}, {{"$limit", list.PageSize}}, {{"$addFields", bson.M{"v": 0}}}}
}

func (list *List) getSortStage(columnsMap map[string]*models.ColumnInfo) (sortStage bson.D) {
	sortFields := bson.M{}
	for _, sortItem := range list.sortFields {
		if sortItem.IsUsed {
			continue
		}
		columnInfo, ok := columnsMap[sortItem.Column]
		if !ok {
			if len(sortFields) == 0 {
				return nil
			}
			sortStage = bson.D{{"$sort", sortFields}}
			return sortStage
		}
		sortFields[columnInfo.Column] = sortItem.Direction
		sortItem.IsUsed = true
	}
	sortStage = bson.D{{"$sort", sortFields}}
	return sortStage
}
func parseSortField(field string) (string, int) {
	if field[:1] == "-" {
		return field[1:], -1
	}
	return field, 1
}
func getSortField(sortItem string, columnsMap map[string]*models.ColumnInfo) (sortField string, direction int) {
	sortField, direction = parseSortField(sortField)
	if columnInfo, ok := columnsMap[sortField]; ok {
		return columnInfo.Column, direction
	}
	return "", 0
}

func (list *List) isPaginatedList() bool {
	return list.ResultType != RESULT_TYPE_EXCEL
}

func (list *List) prepareResult(cursor *mongo.Cursor, context context.Context, dataSliceType interface{}) {
	switch list.ResultType {
	case RESULT_TYPE_PAGINATED_JSON:
		list.prepareJsonResult(cursor, context, dataSliceType)
	case RESULT_TYPE_EXCEL:
		list.prepareExcelResult(cursor, context, dataSliceType)
	}
}

func (list *List) prepareJsonResult(cursor *mongo.Cursor, context context.Context, dataSliceType interface{}) {
	if err := cursor.All(context, dataSliceType); err != nil {
		fmt.Println(err.Error())
		return
	}
	list.Data = dataSliceType
}

// addChangeableResultCountStagesAndFilters if stage has changableResultCount flag then it add to prepared pipeline
// and after each state should add filter stage if filter parameter matched with ResultColumn
// @return pipelines, pipelineIndex, remainingFilters
func (list *List) addChangeableResultCountStagesAndFilters(remainingFilters map[string]string) ([]bson.D, int, map[string]string) {
	pipelineIndex := 0
	var pipelines []bson.D
	filterStage := bson.D{}
	/// add pipelines before sort and pagination stages that maybe change result count
	for pipelineIndex < len(list.dbPipelines) && pipelinesHasChangeableResultCount(list.dbPipelines[pipelineIndex:]) {
		pipeline := list.dbPipelines[pipelineIndex]
		pipelines = append(pipelines, pipeline.Pipeline)
		if len(remainingFilters) != 0 {
			filterStage, remainingFilters = list.generateFilterStage(remainingFilters, pipeline.ResultColumns)
			if filterStage != nil {
				pipelines = append(pipelines, filterStage)
			}
		}
		pipelineIndex++
	}
	return pipelines, pipelineIndex, remainingFilters
}
